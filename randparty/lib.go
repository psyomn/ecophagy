package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/psyomn/ecophagy/common"
)

const (
	AppName = "randparty"

	configFileName = "config.txt"

	defaultConfig = `nouns.txt https://raw.githubusercontent.com/taikuukaits/SimpleWordlists/master/Wordlist-Nouns-All.txt
verbs.txt https://raw.githubusercontent.com/taikuukaits/SimpleWordlists/master/Wordlist-Verbs-All.txt
adjct.txt https://raw.githubusercontent.com/taikuukaits/SimpleWordlists/master/Wordlist-Adjectives-All.txt
`
)

func initialize() { rand.Seed(time.Now().Unix()) }

func dataPath() string { return path.Join(common.DataPath(), AppName) }

func configFile() string { return path.Join(dataPath(), configFileName) }

func createDefaultConfig() error {
	err := os.MkdirAll(dataPath(), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile(), []byte(defaultConfig), 0644)
}

func main() {
	initialize()

	config, err := common.ReadSpaceConfig(configFile())
	if err != nil {
		log.Println("error:", err)

		fmt.Println("you need to supply a config file")
		fmt.Println("each line in the config file is a filename")
		fmt.Println("associated to a url. example: ")
		fmt.Println("")
		fmt.Println("nouns.txt www.someplace.com/nouns.txt")
		fmt.Println("verbs.txt www.someplace.com/verbs.txt")
		fmt.Println("adjct.txt www.someplace.com/adjct.txt")
		fmt.Println("")

		fmt.Print("Would you like to create a default config? [y/n]: ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if text == "n" || text == "no" {
			return
		}

		err = createDefaultConfig()
		if err != nil {
			panic(err)
		}

		log.Println("should download appropriate files next run")

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
