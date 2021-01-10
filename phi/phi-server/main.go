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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/psyomn/ecophagy/common"
	"github.com/psyomn/ecophagy/img"
	"github.com/psyomn/ecophagy/phi/config"
	"github.com/psyomn/ecophagy/phi/phi-server/static"
)

const (
	minUsernameLength = 8
	minPasswordLength = 8

	version = "1.0.0"
)

type errorResponse struct {
	Error string `json:"error"`
}

type session struct {
	configPath string
	config     *config.Config
}

type controller struct {
	server *server
}

func respondWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Println(err)

	errRespJSON, err := json.Marshal(&errorResponse{
		Error: err.Error(),
	})
	if err != nil {
		return
	}
	w.Write(errRespJSON)
}

func (s *controller) handleStatus(w http.ResponseWriter, r *http.Request) {
	type status struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}

	ret := status{Status: "ok", Version: version}
	data, err := json.Marshal(&ret)
	if err != nil {
		log.Println("could not encode message: ", err)
		fmt.Fprintf(w, "Error")
		return
	}

	fmt.Fprintf(w, "%s", data)
}

var mutex sync.Mutex

// POST
func (s *controller) handleRegister(w http.ResponseWriter, r *http.Request) {
	type register struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var regReq register
	error := json.NewDecoder(r.Body).Decode(&regReq)
	if error != nil {
		w.WriteHeader(400)
		log.Println("problem parsing registration request")
		errorResponse := errorResponse{Error: "problem parsing registration request"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)

		return
	}

	if len(regReq.Password) < minPasswordLength {
		w.WriteHeader(400)
		log.Println("problem registering user with small password")
		errorResponse := errorResponse{Error: "passwords must be larger than 8 characters"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)

		return
	}

	if len(regReq.Username) < minUsernameLength {
		w.WriteHeader(400)
		log.Println("problem registering user with small username")
		errorResponse := errorResponse{Error: "problem registering user with small username"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)

		return
	}

	registerError := s.server.registerUser(regReq.Username, regReq.Password, &mutex)
	if registerError != nil {
		w.WriteHeader(400)
		errorResponse := errorResponse{Error: registerError.Error()}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}
		w.Write(errRespJSON)
	}
}

func (s *controller) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq loginRequest
	error := json.NewDecoder(r.Body).Decode(&loginReq)
	if error != nil {
		respondWithError(w, error)
		return
	}

	if err := validateUsername(loginReq.Username); err != nil {
		respondWithError(w, err)
		return
	}

	if err := validatePassword(loginReq.Password); err != nil {
		respondWithError(w, err)
		return
	}

	token, err := s.server.login(loginReq.Username, loginReq.Password)
	if err != nil {
		respondWithError(w, err)
		return
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}

	tokenJSON, err := json.Marshal(&tokenResponse{token})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(tokenJSON)
}

func (s *controller) handleUpload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*150)

	if r.Method != "POST" {
		err := errors.New("only post supported.")
		respondWithError(w, err)
		return
	}

	uriParts := strings.Split(r.RequestURI, "/")
	filename := "default-filename"
	timestamp := "default-timestamp"
	if len(uriParts) == 4 {
		filename = uriParts[2]
		timestamp = uriParts[3]
	} else {
		respondWithError(w, errors.New("url: urls is malformed: "+r.RequestURI))
		return
	}

	parts := strings.Split(r.Header["Authorization"][0], " ")
	if len(parts) != 2 {
		respondWithError(w, errors.New("badauth: expected 'Authorization: token XXX' format"))
		return
	}
	token := parts[1]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, errors.New("could not read body"))
		return
	}

	if username, ok := s.server.session[token]; ok {
		s.server.upload(filename, username, timestamp, body[:])
		return
	} else {
		fmt.Println(username)
		respondWithError(w, errors.New("please login first"))
	}
}

func (s *controller) handleBrowse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(static.DebugPage))
}

func main() {
	if !img.HasExifTool() {
		log.Println("warning: no exif tool found; no support for comment tagging photos")
	}

	fls := &session{config: &config.Config{}}
	flag.StringVar(&fls.configPath, "config", fls.configPath, "path to config")
	flag.Parse()

	{
		bytes, err := common.FileToBytes(fls.configPath)
		if err != nil {
			fmt.Println("error opening file: ", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(bytes, fls.config); err != nil {
			fmt.Println("error unmarshalling: ", err)
		}
	}

	httpServer := http.NewServeMux()
	backend, err := ServerNew()
	if err != nil {
		panic(err)
	}
	controller := controller{backend}

	// REST API
	httpServer.HandleFunc("/status", controller.handleStatus)
	httpServer.HandleFunc("/register", controller.handleRegister)
	httpServer.HandleFunc("/login", controller.handleLogin)
	httpServer.HandleFunc("/upload/", controller.handleUpload)

	// Browser
	httpServer.HandleFunc("/browse", controller.handleBrowse)

	addr := fmt.Sprintf("127.0.0.1:%s", fls.config.Port)
	err = http.ListenAndServe(addr, httpServer)
	if err != nil {
		log.Fatal(err)
	}
}
