package sqlite

import (
	"context"
	"daas_api/pkg/logger"
	"database/sql"
	"encoding/json"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/mattn/go-sqlite3"
)

type CloseFunc func() error

type SQLite struct {
	logger logger.Logger
	db     *sql.DB
	ctx context.Context
	tableName string
}

func CreateSQLite(logger logger.Logger, ctx context.Context, tableName string) (*SQLite, CloseFunc, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		logger.Errorln("Failed to open SQLite DB")
		return nil, nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	// Check if table exists
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, tableName)
	if err != nil {
		logger.Errorw("Error querying sqlite for table",
			"tableName", tableName,
		)
	}
	if !rows.Next() {
		logger.Errorw("Table does not exist",
			"tableName", tableName,
		)
		_, err = db.Exec(`CREATE TABLE ` + tableName + ` (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			phrase TEXT NOT NULL
		)`)
		if err != nil {
			logger.Errorw("Failed to create table",
				"tableName", tableName,
				"err", err,
			)
			return nil, nil, err
		}
	}

	return &SQLite{
		logger: logger,
		db:     db,
		ctx: ctx,
		tableName: tableName,
	}, db.Close, nil
}

func (s *SQLite) Insert(key string, jsonValue map[string]interface{}) error {
	byteData, err := json.Marshal(jsonValue)
	if err != nil {
		s.logger.Errorw("Failed to marshall json",
			"err", err,
			"json", jsonValue,
		)
		return err
	}
	_, err = s.db.Exec("INSERT INTO " + s.tableName + " (key, phrase) values(?, ?)", key, string(byteData))
	if err != nil {
		s.logger.Errorw("Failed to insert phrase",
			"err", err,
			"key", key,
			"phrase", byteData,
		)
		return err
	}
	return nil
}

func (s *SQLite) Get(key string) (map[string]interface{}, error) {
	var id sql.NullInt64
	var returnedKey sql.NullString
	var returnedJson sql.NullString
	err := s.db.QueryRow(`SELECT id, key, phrase FROM phrases WHERE key = ?`, key).Scan(&id, &returnedKey, &returnedJson)
	if err != nil {
		s.logger.Debugw("Failed to query database",
			"err", err,
			"key", key,
		)
		return nil, err
	}

	if !returnedJson.Valid {
		return nil, nil
	} else {
		// Decode JSON content into a struct
		var decodedJson map[string]interface{}
		err = json.Unmarshal([]byte(returnedJson.String), &decodedJson)
		if err != nil {
			s.logger.Errorw("Failed to unmarshal json",
				"err", err,
			)
		}
		return decodedJson, nil
	}
}

func (s *SQLite) GetAll() ([]map[string]interface{}, error) {
	var phrases []map[string]interface{}

	rows, err := s.db.Query("SELECT * FROM "+ s.tableName)
	for rows.Next() {
		var id sql.NullInt64
		var returnedKey sql.NullString
		var returnedJson sql.NullString
		err = rows.Scan(&id, &returnedKey, &returnedJson)
		if err != nil {
			s.logger.Warnw("Skipping entry",
				"err", err,
			)
			continue
		}

		if !returnedJson.Valid {
			s.logger.Warnw("Not valid entry",
				"returnedJson", returnedJson,
			)
			continue
		} else {
			// Decode JSON content into a struct
			var decodedJson map[string]interface{}
			err = json.Unmarshal([]byte(returnedJson.String), &decodedJson)
			if err != nil {
				s.logger.Errorw("Failed to unmarshal json",
					"err", err,
				)
			}
			phrases = append(phrases, decodedJson)
		}
	}
	return phrases, nil
}

func (s *SQLite) Delete(key string) error {
	_, err := s.db.Exec("DELETE FROM " + s.tableName + " WHERE key = ?", key)
	if err != nil {
			s.logger.Errorw("Failed to delete from db",
				"err", err,
				"key", key,
			)
		return err
	}
	return nil
}
