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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/psyomn/ecophagy/common"
	"github.com/psyomn/ecophagy/img"
	"github.com/psyomn/ecophagy/phi/config"
)

const (
	minUsernameLength = 8
	minPasswordLength = 8

	version = "1.0.0"
)

type session struct {
	configPath string
	config     *config.Config
}

func main() {
	common.DefaultLogSettings()

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
			os.Exit(1)
		}

		if err := os.MkdirAll(path.Dir(fls.config.ImagesPath), 0755); err != nil {
			fmt.Println("error creating dir: ", err)
		}

		if err := os.MkdirAll(path.Dir(fls.config.DatabasePath), 0755); err != nil {
			fmt.Println("error creating dir: ", err)
		}
	}

	httpServer := http.NewServeMux()
	backend, err := BackendNew(fls.config.DatabasePath, fls.config.ImagesPath)
	if err != nil {
		panic(err)
	}
	controller := controller{backend}

	// REST API
	httpServer.HandleFunc("/status", controller.handleStatus)
	httpServer.HandleFunc("/register", controller.handleRegister)
	httpServer.HandleFunc("/login", controller.handleLogin)
	httpServer.HandleFunc("/upload/", controller.handleUpload)
	httpServer.HandleFunc("/view/", controller.handleView)

	// Browser
	httpServer.HandleFunc("/browse", controller.handleBrowse)

	addr := fmt.Sprintf("127.0.0.1:%s", fls.config.Port)
	err = http.ListenAndServe(addr, httpServer)
	if err != nil {
		log.Fatal(err)
	}
}
