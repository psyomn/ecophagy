package models

import (
	"database/sql"
	"errors"
)

type Handle struct {
	db *sql.DB
}

var Handler *Handle

func HandleNew(path string) (*Handle, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Handle{db}, nil
}

func (s *Handle) GetRaw() *sql.DB { return s.db }
func (s *Handle) Cleanup()        { s.db.Close() }

func (s *Handle) Execute(stmtstr string, values ...any) (*sql.Result, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(stmtstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)

	if err := tx.Commit(); err != nil {
		errRoll := tx.Rollback()
		return nil, errors.Join(err, errRoll)
	}

	return &result, err
}

func (s *Handle) QueryRow(stmtstr string, values ...any) (*sql.Row, error) {
	stmt, err := s.db.Prepare(stmtstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(values...)

	return row, nil
}

// NB: caller responsible to close rows if no error
func (s *Handle) Query(stmtstr string, values ...any) (*sql.Rows, error) {
	stmt, err := s.db.Prepare(stmtstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	//nolint:sqlclosecheck // the caller is responsible to close the rows
	rows, err := stmt.Query(values...)
	return rows, err
}
