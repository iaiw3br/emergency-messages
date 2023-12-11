package tests

import (
	"errors"
	"fmt"
	"github.com/emergency-messages/pkg/client/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "test_db"
)

var (
	dbURL = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbUser, dbPassword, dbName)
	db    *bun.DB
)

// SetupTestDatabase return test database connection
func SetupTestDatabase() (*bun.DB, error) {
	if db == nil {
		return getDatabase()
	}
	return db, nil
}

// getDatabase get test database connection
func getDatabase() (*bun.DB, error) {
	db = postgres.Connect(dbURL)

	if err := migrateDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// migrateDB drop all tables and create tables for testing
func migrateDB() error {
	f, _ := os.Getwd()
	sourceURL := path.Join(filepath.Dir(f), "migration", "test_migration")
	databaseURL := dbURL + "?sslmode=disable"
	m, err := migrate.New(fmt.Sprintf("file:%s", sourceURL), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	// drop tables
	if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	// create tables
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
