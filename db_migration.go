package main

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pkg/errors"
)

func runMigrationsSource(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "couldn't create psql driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://resources/migrations",
		"postgres", driver)
	if err != nil {
		return errors.Wrap(err, "couldn't create migrate instance")
	}
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return errors.Wrap(err, "couldn't run migrations")
	}
	return nil
}
