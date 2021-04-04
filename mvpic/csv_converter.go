package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

// MakeDbFromCSV assumes that you have never run this before, creates
// tables and populates them properly
func MakeDBFromCSV(dbpath, csvpath, csvFilename string) error {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	// create tables if not exist
	setupQueries := [...]string{
		createMovieMetadataTable,
		createMovieIDIndex,
		createMovieNameIndex,
		createUserRatingsTable,
	}

	for _, query := range setupQueries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	// ingest data
	csvFile, err := os.Open(path.Join(csvpath, csvFilename))
	if err != nil {
		return err
	}

	fmt.Println("reading: ", csvFile)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var movieRec movieRecord

	if _, err := db.Exec("BEGIN TRANSACTION;"); err != nil {
		return err
	}

	for {
		// TODO: skip first line, those are the headers
		fmt.Print(".")
		line, err := reader.Read()

		// nolint
		if errors.Is(err, io.EOF) {
			break // done
		} else if errors.Is(err, csv.ErrTrailingComma) ||
			errors.Is(err, csv.ErrBareQuote) ||
			errors.Is(err, csv.ErrQuote) ||
			errors.Is(err, csv.ErrFieldCount) {
			fmt.Println("Skipping malformed csv line...")
			continue
		} else if err != nil {
			return err
		}

		movieRec.convertLineParts(line)

		stmt, _ := db.Prepare(insertMovie)
		defer stmt.Close()

		if _, err := stmt.Exec(
			movieRec.adult,
			movieRec.belongsToCollection,
			movieRec.budget,
			movieRec.genres,
			movieRec.homepage,
			movieRec.id,
			movieRec.imdbID,
			movieRec.originalLanguage,
			movieRec.originalTitle,
			movieRec.overview,
			movieRec.popularity,
			movieRec.posterPath,
			movieRec.productionCompanies,
			movieRec.productionCountries,
			movieRec.releaseDate,
			movieRec.revenue,
			movieRec.runtime,
			movieRec.spokenLanguages,
			movieRec.status,
			movieRec.tagline,
			movieRec.title,
			movieRec.video,
			movieRec.voteAverage,
			movieRec.voteCount); err != nil {
			return err
		}
	}

	_, err = db.Exec("COMMIT TRANSACTION;")
	return err
}
