package api

import (
	"context"
	"daas/internal/logger"
	"daas/internal/phrase"
	"net/http"

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
    return &Server{
		logger: logger,
		ctx: ctx,
        Router: gin.Default(),
        srv:    &http.Server{Addr: addr, Handler: nil},
		pdb: pdb,
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
