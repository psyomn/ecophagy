package storage

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

const Name = `notes.sqlite3`

var _ embed.FS

//go:embed migrations/*.sql
var schema embed.FS

func MaybeCreateDB(db *sql.DB) error {
	goose.SetBaseFS(schema)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
