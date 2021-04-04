package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
)

func search(t string) []movieRecord {
	// purely search db for entries, and return list of values

	lowerTitle := strings.ToLower(t)
	title := fmt.Sprintf("%%%v%%", lowerTitle) // necessary because can't put parenthesis inside the statement

	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query(searchMovieLike, title)
	if err != nil {
		log.Println(err, "title:", title)
		return nil
	}
	defer rows.Close()

	var ret []movieRecord
	for rows.Next() {
		var mr movieRecord
		err := rows.Scan(
			&mr.adult,
			&mr.belongsToCollection,
			&mr.budget,
			&mr.genres,
			&mr.homepage,
			&mr.id,
			&mr.imdbID,
			&mr.originalLanguage,
			&mr.originalTitle,
			&mr.overview,
			&mr.popularity,
			&mr.posterPath,
			&mr.productionCompanies,
			&mr.productionCountries,
			&mr.releaseDateTime,
			&mr.revenue,
			&mr.runtime,
			&mr.spokenLanguages,
			&mr.status,
			&mr.tagline,
			&mr.title,
			&mr.video,
			&mr.voteAverage,
			&mr.voteCount,
		)

		mr.releaseDate = mr.releaseDateTime.Unix()

		if err != nil {
			log.Println(err)
			continue
		}

		ret = append(ret, mr)
	}

	if rows.Err() != nil {
		log.Println(rows.Err())
	}

	return ret
}

func cliSearch(title string) {
	// presentation logic for cli

	movies := search(title)

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for _, movie := range movies {
		minlen := 0
		if 50 > len(movie.overview) {
			minlen = len(movie.overview)
		} else {
			minlen = 50
		}

		fmt.Fprintf(
			writer,
			"%v\t%v\t%v\t%v\n",
			movie.id, movie.title, movie.releaseDateTime.Year(), movie.overview[0:minlen])
	}
	writer.Flush()
}
