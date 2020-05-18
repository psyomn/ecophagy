package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/psyomn/ecophagy/common"

	_ "github.com/mattn/go-sqlite3"
)

type ratingInfo struct {
	score   int
	comment string
	movieID int
}

func (s *ratingInfo) validate() error {
	if s.score < 1 {
		return errors.New("score can't be below 1")
	}

	if s.score > 10 {
		return errors.New("score can't be above 10 (though some movies deserve that)")
	}

	return nil
}

func cliRate(id string) {
	var ri ratingInfo

	movieID, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	ri.movieID = movieID

	if common.HasEditor() {
		rateThroughEditor(&ri)
		return
	} else {
		rateThroughStdin(&ri)
	}

	rate(&ri)
}

// editor mode: 1st line is the rating of the movie. So you should be
// expected to have a format in the following manner:
//
// 8
// this was a very interesting movie wow
//
func rateThroughEditor(ri *ratingInfo) {
	const (
		sample = `10
first line is your rating. you can leave any comments below.
`
		pattern = AppName
	)

	f, err := ioutil.TempFile(dataPath(), AppName)
	if err != nil {
		panic(err)
	}
	// defer os.Remove(f.Name())

	err = ioutil.WriteFile(f.Name(), []byte(sample), 0666)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(common.Editor(), f.Name())
	fmt.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	maybeContents, err := ioutil.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(maybeContents), common.Newline())

	if len(lines) < 2 {
		panic("movie review not in proper format")
	}

	score, err := strconv.Atoi(lines[0])
	if err != nil {
		panic(err)
	}

	ri.score = score
	ri.comment = strings.Join(lines[1:len(lines)-1], common.Newline())

	rate(ri)
}

// stdin input mode
func rateThroughStdin(ri *ratingInfo) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("score [1-10]: ")
		scoreStr, err := reader.ReadString('\n')
		ri.score, err = strconv.Atoi(scoreStr[0 : len(scoreStr)-1])
		err = ri.validate()

		if err != nil {
			fmt.Println(err)
			continue
		}

		break
	}

	fmt.Println("rating comment: ")
	comment, err := reader.ReadString('\n')
	ri.comment = comment[0 : len(comment)-1]

	if err != nil {
		panic(err)
	}
}

func rate(ri *ratingInfo) {
	db, err := sql.Open("sqlite3", dbPath())

	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(insertUserRating)
	if err != nil {
		panic(err)
	}

	stmt.Exec(ri.movieID, ri.comment, ri.score)
	stmt.Close()

	fmt.Println("rated movie with id:", ri.movieID)
}
