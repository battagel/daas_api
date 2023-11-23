package main

import (
	"context"
	"daas/internal/api"
	"daas/internal/config"
	"daas/internal/db"
	"daas/internal/logger"
	"daas/pkg/redis"
	"fmt"
	"os"
	"os/signal"
)

func run() int {
	// ### Initialisation and Configuration ###
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config", err)
		return 1
	}
	// cli.PrintSplash(config.LoggerEncoding)
	logger, err := logger.CreateZapLogger(config.LoggerLevel, config.LoggerEncoding)
	if err != nil {
		fmt.Println("Error creating logger", err)
		return 1
	}

	// ### Start Program ###
	// Errors can be returned without channels
	logger.Infoln("Starting DaaS")
	logger.Infow("Config",
		"config", config,
	)

	// # Context for async operations #
	// Won't be cancelled until doneChan signal
	asyncCtx, asyncCancel := context.WithCancel(context.Background())
	defer asyncCancel()

	// # Phrase Database #
	// Currently working with Redis backend
	redis, err := redis.CreateRedis(logger, asyncCtx, config.RedisAddress, config.RedisPassword)
	if err != nil {
		logger.Errorw("Error creating Redis Database", err)
		return 1
	}
	pdb, err := db.CreatePhraseDatabase(logger, redis)
	if err != nil {
		logger.Errorw("Error creating Phrase Database", err)
		return 1
	}

	// # API Server #
	// Go gin
	server, err := api.CreateAPIServer(logger, asyncCtx, config.APIAddress, pdb)
	if err != nil {
		logger.Errorw("Error creating API server", err)
		return 1
	}
	server.InitializeRoutes()

	// ### Main Loop ###
	errChan := make(chan error) // Wait for errors
	doneChan := make(chan struct{}) // Unused till shutdown
	shutdownChan := make(chan os.Signal) // To initiate shutdown
	// go events.Start()
	// go metrics.Start()
	go server.Start(doneChan) // Start the API server in goroutine

	// ### Shutdown Sequence ###
	// 1. Stop incoming actions with shutdownChan signal
	// 2. Wait for currently processing actions (API Server and Database) to complete and send doneChan signal
	// 3. Complete remaining shutdown
	signal.Notify(shutdownChan, os.Interrupt)
	select {
	case err = <-errChan:
		logger.Errorw("Error in goroutines",
			"error", err,
		)
		asyncCancel()
	case <-shutdownChan:
		logger.Infoln("Initiating graceful shutdown")
		asyncCancel()
	}
	<-doneChan
	logger.Infoln("Successfully stopped. Exiting")
	return 0
}

func main() {
	os.Exit(run())
}
