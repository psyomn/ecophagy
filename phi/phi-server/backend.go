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
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"

	"github.com/psyomn/ecophagy/common"
	"github.com/psyomn/ecophagy/img"
)

type backend struct {
	session map[string]string
	dbPath  string
	imgPath string
	db      *sql.DB
	mutex   sync.Mutex
}

func BackendNew(dbPath, imgPath string) (*backend, error) {
	var state backend

	state.session = make(map[string]string)
	state.dbPath = dbPath
	state.imgPath = imgPath

	createDbIfNotExist(state.dbPath)
	db, err := getDb(state.dbPath)
	if err != nil {
		return nil, err
	}
	state.db = db
	return &state, nil
}

func getDb(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}

func createDbIfNotExist(dbPath string) {
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return
	}

	db, err := getDb(dbPath)
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

func (s *backend) registerUser(username, password string, mutex *sync.Mutex) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashedPassword, saltStr := common.EncryptPassword(password)

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

func (s *backend) login(username, password string) (string, error) {
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
	if _, err := rand.Read(userToken); err != nil {
		panic(err)
	}

	tokenHex := fmt.Sprintf("%x", userToken)
	s.session[tokenHex] = dbUsername

	log.Println("user", dbUsername, "logged in")
	return tokenHex, nil
}

// TODO: data might be better as a stream type
func (s *backend) upload(filename, username, timestamp string, data []byte) error {
	tm, err := strconv.Atoi(timestamp)
	if err != nil {
		log.Println("warning: could not parse timestamp:", timestamp, ":", err)
	}
	ut := time.Unix(int64(tm), 0)
	imgDir := fmt.Sprintf("%d-%02d-%02d", ut.Year(), ut.Month(), ut.Day())

	imgPath := path.Join(s.imgPath, imgDir, filename)
	if err := os.MkdirAll(path.Dir(imgPath), 0755); err != nil {
		log.Println("could not create img date dir", err)
	}

	fh, err := os.Create(imgPath)
	if err != nil {
		log.Println("could not open file" + err.Error())
		return err
	}

	if _, err := fh.Write(data); err != nil {
		log.Println("could not write file: " + err.Error())
		return err
	}

	if err := fh.Sync(); err != nil {
		return err
	}

	if err := fh.Close(); err != nil {
		return err
	}

	log.Println("upload received", imgPath)

	if img.HasExifTool() {
		var cmt userComment
		cmt.Phi.Username = username
		cmt.Phi.Timestamp = int64(tm)

		// default tags
		cmt.Phi.Tags = []string{"phi", username}
		return img.SetExifComment(filename, string(cmt.toJSON()))
	}

	return nil
}
