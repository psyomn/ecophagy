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
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"git.sr.ht/~psyomn/ecophagy/common"
	"git.sr.ht/~psyomn/ecophagy/img"

	"golang.org/x/crypto/bcrypt"
)

type Backend struct {
	session map[string]string
	dbPath  string
	imgPath string
	db      *sql.DB
	mutex   sync.Mutex
}

func BackendNew(dbPath, imgPath string) (*Backend, error) {
	var state Backend

	state.session = make(map[string]string)
	state.dbPath = dbPath
	state.imgPath = imgPath

	createDBIfNotExist(state.dbPath)
	db, err := getDB(state.dbPath)
	if err != nil {
		return nil, err
	}
	state.db = db
	return &state, nil
}

func getDB(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}

func createDBIfNotExist(dbPath string) {
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return
	}

	db, err := getDB(dbPath)
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

func (s *Backend) registerUser(username, password string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hashedPassword, saltStr := common.EncryptPassword(password)

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	stmt, err := tx.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, hashedPassword, saltStr)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrUsernameTaken
		}

		return err
	}

	return tx.Commit()
}

func (s *Backend) login(username, password string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("could not start transaction: %w", err)
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

	err = stmt.QueryRow(username).Scan(
		&dbUsername, &dbPassword, &dbSalt)

	if err != nil {
		return "", ErrGenericLoginError
	}

	passwordsMatchErr := bcrypt.CompareHashAndPassword(
		[]byte(dbPassword), []byte(dbSalt+password))

	if passwordsMatchErr != nil {
		return "", ErrGenericLoginError
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
func (s *Backend) upload(filename, username, timestamp string, data []byte) error {
	tm, err := strconv.Atoi(timestamp)
	if err != nil {
		log.Println("warning: could not parse timestamp:", timestamp, ":", err)
	}
	ut := time.Unix(int64(tm), 0)
	imgDir := fmt.Sprintf("%d-%02d-%02d", ut.Year(), ut.Month(), ut.Day())

	imgPath := path.Join(s.imgPath, username, imgDir, filename)
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
		return fmt.Errorf("error syncing uploaded file: %w", err)
	}

	if err := fh.Close(); err != nil {
		return fmt.Errorf("error closing uploaded file: %w", err)
	}

	log.Println("upload received", imgPath)

	if img.HasExifTool() {
		var cmt userComment
		cmt.Phi.Username = username
		cmt.Phi.Timestamp = int64(tm)
		cmt.Phi.Tags = []string{"phi", username}

		if err := img.SetExifComment(imgPath, string(cmt.toJSON())); err != nil {
			return fmt.Errorf("error executing exif tool: %w", err)
		}
	}

	return nil
}

func (s *Backend) getImageTags(path string) ([]byte, error) {
	str, err := img.GetExifComment(path)
	return []byte(str), err
}
