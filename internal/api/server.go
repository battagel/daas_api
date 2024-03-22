package api

import (
	"context"
	"daas_api/internal/phrase"
	"daas_api/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type PhraseDatabase interface {
	AddPhrase(phrase.Phrase) error
	GetPhrase(string) (phrase.Phrase, error)
	GetAllPhrases() ([]phrase.Phrase, error)
	DeletePhrase(string) error
}

type Server struct {
	logger   logger.Logger
	ctx      context.Context
	Router   *gin.Engine
	srv      *http.Server
	pdb      PhraseDatabase
	certFile string
	keyFile  string
}

func CreateAPIServer(logger logger.Logger, ctx context.Context, mode, addr, certFile, keyFile string, pdb PhraseDatabase) (*Server, error) {
	// Create a new Gin router
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Apply CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://stackoverflow.com", "http://confluence.eng.nimblestorage.com", "https://rndwiki-pro.its.hpecorp.net"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true // Allow credentials in cross-origin requests
	router.Use(cors.New(config))

	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		logger:   logger,
		ctx:      ctx,
		Router:   router,
		srv:      srv,
		pdb:      pdb,
		certFile: certFile,
		keyFile:  keyFile,
	}, nil
}

func (s *Server) Start() {
	if err := s.srv.ListenAndServeTLS(s.certFile, s.keyFile); err != http.ErrServerClosed {
		s.logger.Errorw("Server error: %v\n", err)
	}
}

func (s *Server) Stop() error {
	err := s.srv.Shutdown(s.ctx)
	if err != nil {
		s.logger.Errorw("Server shutdown error: %v\n", err)
		return err
	}
	s.logger.Debugln("Server gracefully shut down")
	return nil
}
