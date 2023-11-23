package simplejsondb

import (
	"daas/internal/logger"
	"encoding/json"
	"errors"
	"os"
)

type SimpleJsonDB struct {
	logger logger.Logger
	data   map[string]map[string]interface{}
}

func CreateSimpleJsonDB(logger logger.Logger, jsonFilePath string) (*SimpleJsonDB, error) {
	db := &SimpleJsonDB{
		logger: logger,
		data:   make(map[string]map[string]interface{}),
	}
	if err := db.loadJsonFile(jsonFilePath); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *SimpleJsonDB) loadJsonFile(jsonFilePath string) error {
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &db.data); err != nil {
		return err
	}
	return nil
}

func (db *SimpleJsonDB) Get(key string) (map[string]interface{}, error) {
	phraseData, found := db.data[key]
	if !found {
		return nil, errors.New("Key not found")
	}

	return phraseData, nil
}

func (db *SimpleJsonDB) Insert(key string, data map[string]interface{}) error {
	db.data[key] = data
	return nil
}

func (db *SimpleJsonDB) Delete(key string) error {
	if _, found := db.data[key]; !found {
		return errors.New("Key not found")
	}

	delete(db.data, key)
	return nil
}

func (db *SimpleJsonDB) GetAll() map[string]map[string]interface{} {
	return db.data
}
