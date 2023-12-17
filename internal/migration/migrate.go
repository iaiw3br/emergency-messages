package migration

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/source/github"
	_ "github.com/jackc/pgx/stdlib"
)

func RunMigrate(url string) error {
	f, _ := os.Getwd()
	databaseURL := url + "?sslmode=disable"
	sourceURL := path.Join(filepath.Dir(f), "migration")
	m, err := migrate.New(fmt.Sprintf("file:%s", sourceURL), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	// create tables
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
