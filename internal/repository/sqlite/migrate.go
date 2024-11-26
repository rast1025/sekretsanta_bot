package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

func MigrateUp(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not initialize migrator driver")

	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "sqlite3", driver)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("could not initialize migrator: %v", err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	return nil
}
