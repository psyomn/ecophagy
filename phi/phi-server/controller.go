/*
Copyright 2019-2021 Simon Symeonidis (psyomn)

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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/psyomn/ecophagy/common"
	"github.com/psyomn/ecophagy/phi/phi-server/static"
)

type controller struct {
	backend *Backend
}

func (s *controller) checkLogin(r *http.Request) (string, error) {
	parts := strings.Split(r.Header["Authorization"][0], " ")
	if len(parts) != 2 {
		return "", ErrBadAuthHeader
	}

	token := parts[1]
	if username, ok := s.backend.session[token]; ok {
		return username, nil
	}

	return "", ErrNeedLogin
}

func (s *controller) handleRegister(w http.ResponseWriter, r *http.Request) {
	type register struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var regReq register
	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("problem parsing registration request")
		errorResponse := errorResponse{Error: "problem parsing registration request"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}

		if _, err := w.Write(errRespJSON); err != nil {
			log.Println(err)
		}

		return
	}

	// TODO: Disable special characters / accept only alphanum

	if len(regReq.Password) < minPasswordLength {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("problem registering user with small password")
		errorResponse := errorResponse{Error: "passwords must be larger than 8 characters"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}

		if _, err := w.Write(errRespJSON); err != nil {
			log.Println(err)
		}

		return
	}

	if len(regReq.Username) < minUsernameLength {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("problem registering user with small username")
		errorResponse := errorResponse{Error: "problem registering user with small username"}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}

		if _, err := w.Write(errRespJSON); err != nil {
			log.Println(err)
		}

		return
	}

	registerError := s.backend.registerUser(regReq.Username, regReq.Password)
	if registerError != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := errorResponse{Error: registerError.Error()}
		errRespJSON, err := json.Marshal(&errorResponse)
		if err != nil {
			return
		}

		if _, err = w.Write(errRespJSON); err != nil {
			log.Println("problem responding", err)
		}
	}
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

func (s *controller) handleLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq loginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		respondWithError(w, err)
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

	token, err := s.backend.login(loginReq.Username, loginReq.Password)
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

	if _, err := w.Write(tokenJSON); err != nil {
		log.Println("problem with login request:", err)
	}
}

func (s *controller) handleUpload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	username, err := s.checkLogin(r)
	if err != nil {
		respondWithError(w, err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*150)

	if r.Method != "POST" {
		respondWithError(w, ErrMethodNotSupported)
		return
	}

	uriParts, err := common.PartsOfURLSafe(r.RequestURI)
	if err != nil {
		log.Println(err)
		respondWithError(w, ErrMalformedURL)
		return
	}

	if len(uriParts) != 3 {
		respondWithError(w, ErrMalformedURL)
		return
	}

	filename := uriParts[1]
	timestamp := uriParts[2]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, ErrBadBody)
		return
	}

	err = s.backend.upload(filename, username, timestamp, body)
	if err != nil {
		log.Println("could not upload:", err)
	}
}

func (s *controller) handleBrowse(w http.ResponseWriter, r *http.Request) {
	// TODO: might be worth getting rid of this
	_, err := w.Write([]byte(static.DebugPage))
	if err != nil {
		log.Println(err)
	}
}

func (s *controller) handleView(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respondWithError(w, ErrMethodNotSupported)
		return
	}

	username, err := s.checkLogin(r)
	if err != nil {
		respondWithError(w, err)
		return
	}

	uriParts, err := common.PartsOfURLSafe(r.RequestURI)
	if err != nil {
		respondWithError(w, err)
		return
	}

	userPath := path.Join(s.backend.imgPath, username)

	dirName := ""
	pictureName := ""
	switch len(uriParts) {
	case 1:
		goto viewDirs // GET /view
	case 2:
		dirName = uriParts[1]
		goto viewFiles // GET /view/yyyy-mm-dd
	case 3:
		dirName = uriParts[1]
		pictureName = uriParts[2]
		goto fetchFile // GET /view/yyyy-mm-dd/filename
	default:
		respondWithError(w, ErrBadEndpoint)
		return
	}

viewDirs:
	{
		type listDirsResponse struct {
			Dirs []string `json:"dirs"`
		}

		if files, err := filepath.Glob(userPath + "/*"); err != nil {
			respondWithError(w, err)
		} else {
			respJSON, err := json.Marshal(&listDirsResponse{Dirs: files})
			if err != nil {
				log.Println(err)
			}

			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(respJSON); err != nil {
				log.Println(err)
			}
		}

		return
	}

viewFiles:
	{
		type listFilesResponse struct {
			Files []string `json:"files"`
		}

		files, err := filepath.Glob(userPath + "/" + dirName + "/*")

		if err != nil {
			respondWithError(w, err)
			return
		}

		respJSON, err := json.Marshal(&listFilesResponse{Files: files})
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(respJSON); err != nil {
			log.Println(err)
		}

		return
	}

fetchFile:
	{
		w.WriteHeader(http.StatusOK)
		pic := path.Join(userPath, dirName, pictureName)

		fh, err := os.Open(pic)
		if err != nil {
			respondWithError(w, err)
			return
		}
		defer fh.Close()

		bytes, err := io.ReadAll(fh)
		if err != nil {
			log.Println(err)
			return
		}

		if _, err := w.Write(bytes); err != nil {
			log.Println(err)
		}

		return
	}
}

func (s *controller) handleTag(w http.ResponseWriter, r *http.Request) {
	username, err := s.checkLogin(r)
	if err != nil {
		respondWithError(w, ErrNeedLogin)
		return
	}

	uriParts, err := common.PartsOfURLSafe(r.RequestURI)
	if err != nil && len(uriParts) != 3 {
		log.Println(err)
		respondWithError(w, ErrMalformedURL)
		return
	}

	dir := uriParts[1]
	filename := uriParts[2]
	imgPath := path.Join(s.backend.imgPath, username, dir, filename)

	switch r.Method {
	case `GET`:
		goto getTags
	case `PATCH`:
		goto patchTags
	default:
		respondWithError(w, ErrMethodNotSupported)
		return
	}

getTags:
	{
		log.Println("get tags")
		raw, err := s.backend.getImageTags(imgPath)
		if err != nil {
			log.Println("respond with error")
			respondWithError(w, err)
			return
		}

		log.Println("write")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(raw); err != nil {
			log.Println(err)
		}
		log.Println("done")
		return
	}

patchTags:
	{
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`patchtags`)); err != nil {
			log.Println(err)
		}
		return
	}
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

	if _, err := w.Write(errRespJSON); err != nil {
		log.Println("problem responding with error:", err)
	}
}
