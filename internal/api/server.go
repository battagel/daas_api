package api

import (
	"context"
	"daas_api/internal/phrase"
	"daas_api/pkg/logger"
	"net/http"

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
	logger logger.Logger
	ctx 	context.Context
    Router *gin.Engine
    srv    *http.Server
	pdb PhraseDatabase
}

func CreateAPIServer(logger logger.Logger, ctx context.Context, addr string, pdb PhraseDatabase) (*Server, error) {
    // Create a new Gin router
    router := gin.Default()

    // Apply CORS middleware
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"https://stackoverflow.com", "http://confluence.eng.nimblestorage.com", "https://rndwiki-pro.its.hpecorp.net"}
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
    config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true // Allow credentials in cross-origin requests
    router.Use(cors.New(config))

    return &Server{
        logger: logger,
        ctx:    ctx,
        Router: router,
        srv:    &http.Server{Addr: addr, Handler: router},
        pdb:    pdb,
    }, nil
}

func (s *Server) Start(doneChan chan struct{}) {
    s.srv.Handler = s.Router
	go func() {
        select {
        case <-s.ctx.Done():
            // The context is canceled, initiate server shutdown
            err := s.srv.Shutdown(context.Background())
			if err != nil {
                s.logger.Errorw("Server shutdown error: %v\n", err)
			}
			s.logger.Debugln("Server gracefully shut down")
			doneChan <- struct{}{} // Notify that we have finished processing
        }
    }()

    if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
        s.logger.Errorw("Server error: %v\n", err)
    }
}
