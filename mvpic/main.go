package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/psyomn/ecophagy/common"
)

const (
	// AppName is the application name under ecophagy
	AppName = `mvpic`

	// advice is where the data set was originally taken from
	advice = `https://www.kaggle.com/rounakbanik/the-movies-dataset`
)

var (
	cmdUpdateDatabase = false
	cmdShowConfigPath = false
	cmdSearch         = ""
	cmdDesc           = ""
)

func dataPath() string {
	return path.Join(common.DataPath(), AppName)
}

func dbPath() string {
	return path.Join(dataPath(), "movies.sqlite3")
}

func expectedFiles() map[string]string {
	return map[string]string{
		`credits.csv`:         `d62a1d7d652d324bebee68782f664fc9`,
		`keywords.csv`:        `7d0f5123e11917fa1bea011448e5f73d`,
		`links.csv`:           `b9b8fe775557e10e35418410499775cf`,
		`links_small.csv`:     `917eddf52079d6ce5c6cfd63b17515f7`,
		`movies_metadata.csv`: `42bf3ef8c208a01a4776955875978b1e`,
		`ratings.csv`:         `f640a181f6fa0b2e3294f786fa350ccc`,
		`ratings_small.csv`:   `8864480f98416ccecaf17aa5263bbea1`,
	}
}

// TODO this can be refactored quite well
func checkFiles() error {
	expectedFiles := expectedFiles()

	files, err := common.FileList(dataPath())
	if err != nil {
		return err
	}

	for _, file := range files {
		justFileName := filepath.Base(file)
		expectedHash, ok := expectedFiles[justFileName]

		if !ok {
			// non-interesting file, skip
			continue
		}

		actualHash, err := common.FileToMd5Sum(file)
		if err != nil {
			return err
		}

		if actualHash != expectedHash {
			errorMsg := fmt.Sprintf("%s bad file hash: expected: %s, got: %s",
				file, expectedHash, actualHash)
			return errors.New(errorMsg)
		}
	}

	return nil
}

func setupFiles() {
	moviePath := dataPath()

	if !common.PathExists(moviePath) {
		fmt.Println("you might want to download data into", moviePath)
		fmt.Println("the data can be found here: ", advice)
		err := os.MkdirAll(moviePath, 0766)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err := checkFiles()
	if err != nil {
		panic(err)
	}

}

func setupDB() {
	if err := MakeDbFromCSV(dbPath(), dataPath(), "movies_metadata.csv"); err != nil {
		panic(err)
	}
}

func setupFlags() {
	flag.BoolVar(&cmdUpdateDatabase, "update", cmdUpdateDatabase, "update database with CSV")
	flag.BoolVar(&cmdShowConfigPath, "show-config-path", cmdShowConfigPath, "print the path to your app data")
	flag.StringVar(&cmdSearch, "search", cmdSearch, "search movie by title")
	flag.StringVar(&cmdDesc, "desc", cmdDesc, "describe movie given id")
	flag.Parse()
}

func main() {
	if !common.HasHome() {
		panic("you need to have a home in order to run this")
	}

	setupFlags()

	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(1)
	}

	if cmdUpdateDatabase {
		setupFiles()
		setupDB()
		return
	}

	if cmdShowConfigPath {
		fmt.Println(dataPath())
		return
	}

	if cmdSearch != "" {
		cliSearch(cmdSearch)
		return
	}

	if cmdDesc != "" {
		cliDesc(cmdDesc)
		return
	}
}
