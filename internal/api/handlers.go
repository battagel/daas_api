package api

import (
	"daas_api/internal/phrase"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetAllPhrases(c *gin.Context) {
	foundPhrases, err := s.pdb.GetAllPhrases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying for all phrases"})
		return
	}

	s.logger.Debugw("Number of phrases",
		"num", len(foundPhrases),
	)
	if len(foundPhrases) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No phrases found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"return": foundPhrases})
}

func (s *Server) CreatePhrase(c *gin.Context) {
	var newPhrase phrase.Phrase
	if err := c.ShouldBindJSON(&newPhrase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phraseExists, err := s.pdb.GetPhrase(newPhrase.Phrase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when querying db for conflict", "key": newPhrase.Phrase})
		return
	}

	if !reflect.DeepEqual(phraseExists, phrase.Phrase{}) {
		c.JSON(http.StatusConflict, gin.H{"error": "Phrase already exists. Maybe try updating instead?"})
		return
	}

	err = s.pdb.AddPhrase(newPhrase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save phrase"})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "Phrase created successfully", "phrase": newPhrase})
}

func (s *Server) GetPhrase(c *gin.Context) {
	// Need to do some html url encoding stuff here
	key := c.Param("key")
	s.logger.Debugw("Key Phrase from URL",
		"key", key,
	)

	foundPhrase, err := s.pdb.GetPhrase(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when querying db for key", "key": key})
		return
	}
	if reflect.DeepEqual(foundPhrase, phrase.Phrase{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phrase does not exist", "key": key})
		return
	}
	s.logger.Debugw("Got phrase from db",
		"phrase", foundPhrase,
	)

	c.JSON(http.StatusOK, gin.H{"return": foundPhrase})
}

func (s *Server) UpdatePhrase(c *gin.Context) {
	key := c.Param("key")
	s.logger.Debugw("Key Phrase from URL",
		"key", key,
	)

	foundPhrase, err := s.pdb.GetPhrase(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when querying db for key", "key": key})
		return
	}
	if reflect.DeepEqual(foundPhrase, phrase.Phrase{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phrase does not exist", "key": key})
		return
	}
	s.logger.Debugw("Found phrase to update db",
		"phrase", foundPhrase,
	)
	//Override phrase by deleting then creating
	err = s.pdb.DeletePhrase(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when trying to delete phrase", "key": key})
		return
	}
	s.logger.Debugw("Successfully deleted old phrase",
		"key", key,
	)

	var newPhrase phrase.Phrase
	if err := c.ShouldBindJSON(&newPhrase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.pdb.AddPhrase(newPhrase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new phrase"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Phrase updated successfully", "phrase": newPhrase})
}

func (s *Server) DeletePhrase(c *gin.Context) {
	key := c.Param("key")
	s.logger.Debugw("Key Phrase from URL",
		"key", key,
	)

	// Check if exists before deleting
	foundPhrase, err := s.pdb.GetPhrase(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying key before deleting", "key": key})
		return
	}
	if reflect.DeepEqual(foundPhrase, phrase.Phrase{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phrase does not exist", "key": key})
		return
	}

	err = s.pdb.DeletePhrase(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when trying to delete phrase", "key": key})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Phrase deleted successfully", "key": key})
}

func (s *Server) BulkImportPhrases(c *gin.Context) {
	// Implement your code to bulk import items into db
	// ...
}
