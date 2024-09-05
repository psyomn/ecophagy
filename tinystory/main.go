package main

import (
	"flag"
	"fmt"
	"os"

	"git.sr.ht/~psyomn/ecophagy/tinystory/lib"
)

func makeFlags(sess *tinystory.Session) {
	flag.StringVar(&sess.Host, "host", sess.Host, "specify host to bind server")
	flag.StringVar(&sess.Port, "port", sess.Port, "specify port to bind server")
	flag.StringVar(&sess.Repository, "repository", sess.Repository, "specify story repository")
	flag.StringVar(&sess.Assets, "assets", sess.Assets, "specify the assets root path")
	flag.BoolVar(&sess.ExperimentalParser, "experimental-parser", sess.ExperimentalParser, "use experimental parser")
	flag.Parse()
}

func main() {
	sess := tinystory.MakeDefaultSession()
	makeFlags(sess)

	var docs []tinystory.Document

	if sess.ExperimentalParser {
		fmt.Println("using experimental parser ...")

		d, err := tinystory.ParseAllInDirExt(sess.Repository)
		if err != nil {
			fmt.Printf("experimental: error parsing stories: %s\n", err.Error())
			os.Exit(1)
		}
		docs = d
	} else {
		d, err := tinystory.ParseAllInDir(sess.Repository)
		if err != nil {
			fmt.Printf("error parsing stories: %s\n", err.Error())
			os.Exit(1)
		}
		docs = d
	}

	server, err := tinystory.ServerNew(sess, docs)
	if err != nil {
		fmt.Println("could not start server:", err)
	}

	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}
