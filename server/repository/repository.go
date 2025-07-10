package repository

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"

	"terra9.it/checkmate/core"
	"terra9.it/checkmate/server/models"
)

func NewStorage() *sqlite3.Storage {
	// Init SQLite3 database
	db, err := sql.Open("sqlite3", "./db/fiber.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Storage package can create this table for you at init time
	// but for the purpose of this example I created it manually
	// expanding its structure with an "u" column to better query
	// all user-related sessions.
	query := `CREATE TABLE IF NOT EXISTS sessions (
		  k  VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT '',
		  v  BLOB NOT NULL,
		  e  BIGINT NOT NULL DEFAULT '0',
		  u  TEXT);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	gob.Register(&core.ProjectExport{})

	// Init sessions store
	storage := sqlite3.New(sqlite3.Config{
		Database:        "./db/fiber.db",
		Table:           "sessions",
		Reset:           false,
		GCInterval:      60 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 60 * time.Second,
	})
	return storage
}

// Simulate a database call
func FindByCredentials(email, password string) (*models.User, error) {
	// Here you would query your database for the user with the given email
	if email == "test@mail.com" && password == "test12345" {
		return &models.User{
			ID:        1,
			Email:     "test@mail.com",
			Password:  "test12345",
			SessionID: utils.UUIDv4(),
		}, nil
	}
	return nil, errors.New("user not found")
}
