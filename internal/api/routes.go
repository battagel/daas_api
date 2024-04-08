package api

// These are all the API routes available. These are then handled by their
// respective functions in handlers.go

func (s *Server) InitializeRoutes() {
	s.Router.GET("/phrase", s.GetAllPhrases)
	s.Router.POST("/phrase/:key", s.CreatePhrase)
	s.Router.GET("/phrase/:key", s.GetPhrase)
	s.Router.PUT("/phrase/:key", s.UpdatePhrase)
	s.Router.DELETE("/phrase/:key", s.DeletePhrase)
	s.Router.POST("/import", s.BulkImportPhrases)
}
