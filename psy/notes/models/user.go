package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	sqlUserCreate   = `INSERT INTO users (name,password,created_at,updated_at) VALUES (?,?,?,?);`
	sqlUserFindName = `SELECT * from users WHERE name = ?;`
)

type User struct {
	ID             int
	Name           string
	Password       string
	HashedPassword []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func UserFindByName(name string) (*User, error) {
	row, err := Handler.QueryRow(sqlUserFindName, name)
	if err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	usr := User{}
	if err := row.Scan(&usr.ID, &usr.Name, &usr.HashedPassword, &usr.CreatedAt, &usr.UpdatedAt); err != nil {
		log.Println(err)
		return nil, err
	}

	return &usr, nil
}

func UserLogin(name, password string) (*User, error) {
	usr, err := UserFindByName(name)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(usr.HashedPassword, []byte(password)); err != nil {
		return nil, fmt.Errorf("%w: by %v", ErrLoginIncorrect, name)
	}

	return usr, nil
}

func UserCreate(name, password string) (*sql.Result, error) {
	penc, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	r, err := Handler.Execute(sqlUserCreate, name, penc, time.Now(), time.Now())
	return r, err
}
