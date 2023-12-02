package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/emergency-messages/pkg/client/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "test_db"
)

var (
	dbURL = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", dbUser, dbPassword, dbName)
	db    *pgx.Conn
)

func SetupTestDatabase() (*pgx.Conn, error) {
	if db == nil {
		return getDatabase()
	}
	return db, nil
}

func getDatabase() (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := postgres.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB() error {
	f, _ := os.Getwd()
	sourceURL := path.Join(filepath.Dir(f), "migration", "test_migration")
	databaseURL := dbURL + "?sslmode=disable"
	m, err := migrate.New(fmt.Sprintf("file:%s", sourceURL), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	// TODO: drop tables after all tests
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
