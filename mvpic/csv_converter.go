package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createMovieMetadataTable = `
CREATE TABLE IF NOT EXISTS movies (
  adult                   BOOLEAN,
  belongs_to_collection   STRING, -- json obj
  budget                  BIGINT,
  genres                  TEXT,
  homepage                TEXT,
  id                      INTEGER PRIMARY KEY,
  imdb_id                 TEXT,
  original_language       TEXT,
  original_title          TEXT,
  overview                TEXT,
  popularity              DOUBLE,
  poster_path             TEXT,
  production_companies    TEXT,
  production_countries    TEXT,
  release_date            DATETIME,
  revenue                 BIGINT,
  runtime                 DOUBLE,
  spoken_languages        TEXT,
  status                  TEXT,
  tagline                 TEXT,
  title                   TEXT,
  video                   BOOLEAN,
  vote_average            DOUBLE,
  vote_count              INT
);`
	insertMovie = `INSERT INTO movies values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	createMovieIdIndex   = `CREATE INDEX IF NOT EXISTS movie_id_index ON movies(id);`
	createMovieNameIndex = `CREATE INDEX IF NOT EXISTS movie_original_title_index ON movies(original_title)`
)

type movieRecord struct {
	adult               bool
	belongsToCollection string
	budget              uint64
	genres              string
	homepage            string
	id                  uint64
	imdbID              string
	originalLanguage    string
	originalTitle       string
	overview            string
	popularity          float64
	posterPath          string
	productionCompanies string
	productionCountries string
	releaseDate         string
	revenue             uint64
	runtime             float64
	spokenLanguages     string
	status              string
	tagline             string
	title               string
	video               bool
	voteAverage         float64
	voteCount           uint64
}

func (s *movieRecord) convertLineParts(line []string) {
	s.adult = line[0] == "True"
	s.belongsToCollection = line[1]

	{
		budget, err := strconv.ParseUint(line[2], 10, 64)
		if err == nil {
			s.budget = budget
		}
	}

	s.genres = line[3]
	s.homepage = line[4]
	{
		id, err := strconv.ParseUint(line[5], 10, 64)
		if err == nil {
			s.id = id
		}
	}

	s.imdbID = line[6]
	s.originalLanguage = line[7]
	s.originalTitle = line[8]
	s.overview = line[9]

	{
		f, err := strconv.ParseFloat(line[10], 64)
		if err == nil {
			s.popularity = f
		}
	}

	s.posterPath = line[11]
	s.productionCompanies = line[12]
	s.productionCountries = line[13]
	s.releaseDate = line[14] // TODO: unix timestamps?

	{
		u, err := strconv.ParseUint(line[15], 10, 64)
		if err == nil {
			s.revenue = u
		}
	}

	{
		f, err := strconv.ParseFloat(line[16], 64)
		if err == nil {
			s.runtime = f
		}
	}

	s.spokenLanguages = line[17]
	s.status = line[18]
	s.tagline = line[19]
	s.title = line[20]
	s.video = line[21] == "True"
	{
		f, err := strconv.ParseFloat(line[22], 64)
		if err == nil {
			s.voteAverage = f
		}
	}

	{
		u, err := strconv.ParseUint(line[23], 10, 64)
		if err == nil {
			s.voteCount = u
		}
	}
}

// MakeDbFromCSV assumes that you have never run this before, creates
// tables and populates them properly
func MakeDbFromCSV(dbpath, csvpath, csvFilename string) error {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	// create tables if not exist
	setupQueries := [...]string{
		createMovieMetadataTable,
		createMovieIdIndex,
		createMovieNameIndex,
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

	for {
		fmt.Print(".")
		line, err := reader.Read()

		if err == io.EOF {
			break // done
		} else if csvErr, ok := err.(*csv.ParseError); ok && csvErr.Err == csv.ErrFieldCount {
			fmt.Println("Skipping malformed csv line...")
			continue
		} else if err != nil {
			return err
		}

		movieRec.convertLineParts(line)

		stmt, _ := db.Prepare(insertMovie)

		stmt.Exec(
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
			movieRec.voteCount)

		stmt.Close()

	}

	return nil
}
