package db

import (
	"daas_api/pkg/logger"
	"daas_api/internal/phrase"
)

type Database interface {
	Insert(string, map[string]interface{}) error
	Get(string) (map[string]interface{}, error)
	GetAll() ([]map[string]interface{}, error)
	Delete(string) error
}

type PhraseDatabase struct {
	logger logger.Logger
	database Database
}

func CreatePhraseDatabase(logger logger.Logger, database Database) (*PhraseDatabase, error) {
	return &PhraseDatabase{
		logger: logger,
		database: database,
	}, nil
}

func (pdb *PhraseDatabase) AddPhrase(phrase phrase.Phrase) error {
	pdb.logger.Debugw("Adding phrase to db",
		"phrase", phrase)
	err := pdb.database.Insert(phrase.Phrase, phrase.ToMap())
	if err != nil {
		return err
	}
	return nil
}

func (pdb *PhraseDatabase) GetPhrase(key string) (phrase.Phrase, error) {
	pdb.logger.Debugw("Key phrase to query",
		"key", key)
	rawData, err := pdb.database.Get(key)
	if err != nil {
		pdb.logger.Errorw("Error getting phrase",
			"error", err,
			"key", key,
		)
		return phrase.Phrase{}, err
	}
	if rawData == nil {
		pdb.logger.Debugw("Phrase not found in db",
			"key", key,
		)
		return phrase.Phrase{}, nil
	}
	var newPhrase phrase.Phrase
	err = newPhrase.ToPhrase(rawData)
	if err != nil {
		pdb.logger.Errorw("Error parsing rawData into new phrase",
			"err", err,
			"rawData", rawData,
		)
		return phrase.Phrase{}, err
	}
	return newPhrase, nil
}

func (pdb *PhraseDatabase) GetAllPhrases() ([]phrase.Phrase, error) {

	rawPhrases, err := pdb.database.GetAll()
	if err != nil {
		return nil, err
	}
	if len(rawPhrases) < 1 {
		pdb.logger.Debugln("No phrases found")
		return nil, nil
	}

	var phrases []phrase.Phrase
	for _, rawPhrase := range rawPhrases {
		pdb.logger.Debugw("Converting phrase",
			"rawPhrase", rawPhrase,
		)
		var phrase phrase.Phrase
		err := phrase.ToPhrase(rawPhrase)
		if err != nil {
			return nil, err
		}
		phrases = append(phrases, phrase)
	}
	return phrases, nil
}

func (pdb *PhraseDatabase) DeletePhrase(key string) error {
	err := pdb.database.Delete(key)
	if err != nil {
		return err
	}

	return nil
}
