package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/psyomn/ecophagy/tinystory/lib"
)

type Session struct {
	host       string
	port       string
	repository string
}

func MakeDefaultSession() *Session {
	return &Session{
		host:       "127.0.0.1",
		port:       "9090",
		repository: "./stories",
		assets:     "./assets",
	}
}

func makeFlags(sess *Session) {
	flag.StringVar(&sess.host, "host", sess.host, "specify host to bind server")
	flag.StringVar(&sess.port, "port", sess.port, "specify port to bind server")
	flag.StringVar(&sess.repository, "repository", sess.repository, "specify story repository")
	flag.Parse()
}

func main() {
	sess := MakeDefaultSession()
	makeFlags(sess)

	docs, err := tinystory.ParseAllInDir(sess.repository)
	if err != nil {
		fmt.Printf("error parsing stories: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(docs)
}
