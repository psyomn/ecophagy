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
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type serverState struct {
	session map[string]string
	mutex   sync.Mutex
}

var (
	dbHandle *sql.DB
	dbMutex  sync.Mutex
	srvState serverState
)

func init() {
	createDbIfNotExist()
	srvState.session = make(map[string]string)
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

func registerUser(username, password string, mutex *sync.Mutex) error {
	mutex.Lock()
	defer mutex.Unlock()

	hashedPassword, saltStr := encryptPassword(password)
	db, err := getDb()
	if err != nil {
		return fmt.Errorf("could not open db to store user: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
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

func login(username, password string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	db, err := getDb()
	if err != nil {
		return "", fmt.Errorf("could not open db to store user: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
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

	srvState.mutex.Lock()
	defer srvState.mutex.Unlock()
	srvState.session[tokenHex] = dbUsername

	return tokenHex, nil
}

func upload(filename, username, timestamp string, data []byte) error {
	// TODO: I wonder if instead of data []byte we should have stream
	//       access instead (io.Reader)
	fh, err := os.Create(filename)
	if err != nil {
		log.Println("could not open file" + err.Error())
		return err
	}
	defer fh.Close()

	dataToWrite := data

	_, err = fh.Write(dataToWrite)
	if err != nil {
		log.Println("could not write file: " + err.Error())
		return err
	}
	fh.Sync()

	return errors.New("not implemented")
}
