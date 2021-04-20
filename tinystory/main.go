package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/psyomn/ecophagy/tinystory/lib"
)

func makeFlags(sess *tinystory.Session) {
	flag.StringVar(&sess.Host, "host", sess.Host, "specify host to bind server")
	flag.StringVar(&sess.Port, "port", sess.Port, "specify port to bind server")
	flag.StringVar(&sess.Repository, "repository", sess.Repository, "specify story repository")
	flag.StringVar(&sess.Assets, "assets", sess.Assets, "specify the assets root path")
	flag.StringVar(&sess.ExperimentalParser, "experimental-parser", sess.ExperimentalParser, "use experimental parser")
	flag.Parse()
}

func main() {
	sess := tinystory.MakeDefaultSession()
	makeFlags(sess)

	if sess.ExperimentalParser != "" {
		// TODO experimental for now
		doc, err := tinystory.ParseTinyStoryFormatFile(sess.ExperimentalParser)
		if err != nil {
			fmt.Println("error parsing tinystory format:", err)
			return
		}

		fmt.Println(doc)

		return
	}

	// TODO: there should be a less bleedy initialization here
	docs, err := tinystory.ParseAllInDir(sess.Repository)
	if err != nil {
		fmt.Printf("error parsing stories: %s\n", err.Error())
		os.Exit(1)
	}

	server, err := tinystory.ServerNew(sess, docs)
	if err != nil {
		fmt.Println("could not start server:", err)
	}

	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}
