package main

import (
	"context"
	"daas_api/internal/api"
	"daas_api/internal/db"
	"daas_api/pkg/config"
	"daas_api/pkg/logger"
	"daas_api/pkg/sqlite"
	"fmt"
	"os"
	"os/signal"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)


func run() exitCode {
	// ### Initialisation and Configuration ###
	config, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config", err)
		return exitError
	}
	// cli.PrintSplash(config.LoggerEncoding)
	logger, err := logger.CreateZapLogger(config.LoggerLevel, config.LoggerEncoding)
	if err != nil {
		fmt.Println("Error creating logger", err)
		return exitError
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
	backend, backendClose, err := sqlite.CreateSQLite(logger, asyncCtx, config.SQLiteTableName)
	if err != nil {
		logger.Errorw("Error creating Sqlite Database", err)
		return exitError
	}
	defer backendClose()
	pdb, err := db.CreatePhraseDatabase(logger, backend)
	if err != nil {
		logger.Errorw("Error creating Phrase Database", err)
		return exitError
	}

	// # API Server #
	server, err := api.CreateAPIServer(logger, asyncCtx, config.APIMode, config.APIAddress, config.APICert, config.APIKey, pdb)
	if err != nil {
		logger.Errorw("Error creating API server", err)
		return exitError
	}
	server.InitializeRoutes()

	// ### Main Loop ###
	shutdownChan := make(chan os.Signal) // To initiate shutdown
	// go events.Start()
	// go metrics.Start()
	go server.Start() // Start the API server in goroutine


	// ### Shutdown Sequence ###
	// 1. Stop incoming actions with shutdownChan signal
	// 2. Wait for currently processing actions (API Server and Database) to complete and send WaitGroup done.
	// 3. Complete remaining shutdown
	signal.Notify(shutdownChan, os.Interrupt)

	<-shutdownChan

	logger.Infoln("Initiating graceful shutdown")
	err = server.Stop()
	if err != nil {
		logger.Errorln("Failed to stop API server")
		return exitError
	}
	logger.Infoln("Successfully stopped. Exiting")
	return exitOK
}

func main() {
	os.Exit(int(run()))
}
