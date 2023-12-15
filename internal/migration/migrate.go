package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"path"
	"path/filepath"
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
