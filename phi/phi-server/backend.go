/*
Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/psyomn/ecophagy/img"
)

type server struct {
	session map[string]string
	db      *sql.DB
	mutex   sync.Mutex
}

func ServerNew() (*server, error) {
	var state server
	state.session = make(map[string]string)

	createDbIfNotExist()
	db, err := getDb() // TODO get rid of global
	if err != nil {
		return nil, err
	}
	state.db = db
	return &state, nil
}

func getDb() (*sql.DB, error) {
	return sql.Open("sqlite3", dbName)
}

func createDbIfNotExist() {
	if _, err := os.Stat(dbName); !os.IsNotExist(err) {
		return
	}

	db, err := getDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err != nil {
		panic(fmt.Sprintf("could not open db: %v", err))
	}

	_, err = db.Exec(phiSchema)
	if err != nil {
		panic(fmt.Sprintf("could not create schema: %v", err))
	}
}

func encryptPassword(password string) (string, string) {
	saltBytes := make([]byte, 8)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}
	saltStr := string(saltBytes)
	saltedPassword := []byte(saltStr + password)

	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, 14)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword), saltStr
}

func encryptPasswordWithSalt(password string, salt string) string {
	saltedPassword := []byte(salt + password)
	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, 14)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func (s *server) registerUser(username, password string, mutex *sync.Mutex) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashedPassword, saltStr := encryptPassword(password)

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	stmt, err := tx.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, hashedPassword, saltStr)
	if err != nil {
		tx.Rollback()

		if strings.Contains(err.Error(), "UNIQUE") {
			return errors.New("username has been taken")
		}

		return err
	}

	tx.Commit()
	return nil
}

func (s *server) login(username, password string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("could not start transaction: %v", err)
	}

	stmt, err := tx.Prepare(loginUserSQL)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	//
	// Check if password matches
	//
	var (
		dbUsername string
		dbPassword string
		dbSalt     string
	)
	const genericLoginError = "username or passwords do not match"

	err = stmt.QueryRow(username).Scan(
		&dbUsername, &dbPassword, &dbSalt)

	if err != nil {
		return "", errors.New(genericLoginError)
	}

	passwordsMatchErr := bcrypt.CompareHashAndPassword(
		[]byte(dbPassword), []byte(dbSalt+password))

	if passwordsMatchErr != nil {
		return "", errors.New(genericLoginError)
	}

	userToken := make([]byte, 32)
	_, err = rand.Read(userToken)
	if err != nil {
		panic(err)
	}
	tokenHex := fmt.Sprintf("%x", userToken)

	s.session[tokenHex] = dbUsername

	return tokenHex, nil
}

func (s *server) upload(filename, username, timestamp string, data []byte) error {
	// TODO: I wonder if instead of data []byte we should have stream
	//       access instead (io.Reader)
	fh, err := os.Create(filename)
	if err != nil {
		log.Println("could not open file" + err.Error())
		return err
	}

	_, err = fh.Write(data)
	if err != nil {
		log.Println("could not write file: " + err.Error())
		return err
	}
	err = fh.Sync()
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		return err
	}

	if img.HasExifTool() {
		var cmt userComment
		cmt.Phi.Username = username
		tm, err := strconv.Atoi(timestamp)
		if err == nil {
			cmt.Phi.Timestamp = int64(tm)
		}
		// default tags
		cmt.Phi.Tags = []string{"phi", username}
		return img.SetExifComment(filename, string(cmt.toJSON()))
	}

	return nil
}
