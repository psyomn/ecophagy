package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/psyomn/ecophagy/common"
)

func initialize() {
	rand.Seed(time.Now().Unix())
}

func main() {
	initialize()

	config, err := readSpaceConfig("config.txt")
	if err != nil {
		fmt.Println("you need to supply a config file")
		fmt.Println("each line in the config file is a filename")
		fmt.Println("associated to a url. example: ")
		fmt.Println("")
		fmt.Println("nouns.txt www.someplace.com/nouns.txt")
		fmt.Println("verbs.txt www.someplace.com/verbs.txt")
		fmt.Println("adjct.txt www.someplace.com/adjct.txt")
		return
	}

	wantedFiles := config

	for k, v := range wantedFiles {
		fmt.Println("Checking/Downloading:", v, "...")
		err := common.DownloadIfNotExist(k, v)

		if err != nil {
			fmt.Println(err)
		}
	}

	nouns, err := common.FileToLines("nouns.txt")
	if err != nil {
		panic(err)
	}

	verbs, err := common.FileToLines("verbs.txt")
	if err != nil {
		panic(err)
	}

	adjct, err := common.FileToLines("adjct.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("==================================")
	fmt.Println("nouns: ", len(nouns))
	fmt.Println("verbs: ", len(verbs))
	fmt.Println("adjct: ", len(adjct))
	fmt.Println("==================================")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(
			"%s %s %s with a %s %s",
			sampleArray(adjct),
			sampleArray(nouns),
			sampleArray(verbs),
			sampleArray(adjct),
			sampleArray(nouns),
		)

		reader.ReadString('\n')
	}
}
