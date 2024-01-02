package sqlite

import (
	"context"
	"daas_api/pkg/logger"
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/golang-migrate/migrate/v4"
)

type SQLite struct {
	logger logger.Logger
	db     *sql.DB
	ctx context.Context
	tableName string
}

func CreateSQLite(logger logger.Logger, ctx context.Context) (*SQLite, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		logger.Errorln("Failed to open SQLite DB")
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &SQLite{
		logger: logger,
		db:     db,
		ctx: ctx,
		tableName: "phrases",
	}, nil
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
	s.logger.Debugw("Unmars")
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
	row, err := s.db.Query(`SELECT id, key, phrase FROM phrases WHERE key = ?`, key)
	if err != nil {
		s.logger.Debugw("Failed to query database",
			"err", err,
			"key", key,
		)
		return nil, err
	}
	defer row.Close()

	var id sql.NullInt64
	var returnedKey sql.NullString
	var returnedJson sql.NullString
	row.Next()
	err = row.Scan(&id, &returnedKey, &returnedJson)
	if err != nil {
		s.logger.Debugw("Phrase not found in query",
			"key", key,
		)
		return nil, nil
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
