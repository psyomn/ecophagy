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

func desc(id string) movieRecord {
	var ret movieRecord

	db, err := sql.Open("sqlite3", dbPath())
	if err != nil {
		log.Println(err)
		return ret
	}
	defer db.Close()

	rows, err := db.Query(descMovie, id)
	if err != nil {
		log.Println(err)
		return ret
	}
	if rows.Err() != nil {
		log.Println(rows.Err())
		return ret
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&ret.adult,
		&ret.belongsToCollection,
		&ret.budget,
		&ret.genres,
		&ret.homepage,
		&ret.id,
		&ret.imdbID,
		&ret.originalLanguage,
		&ret.originalTitle,
		&ret.overview,
		&ret.popularity,
		&ret.posterPath,
		&ret.productionCompanies,
		&ret.productionCountries,
		&ret.releaseDateTime,
		&ret.revenue,
		&ret.runtime,
		&ret.spokenLanguages,
		&ret.status,
		&ret.tagline,
		&ret.title,
		&ret.video,
		&ret.voteAverage,
		&ret.voteCount,
	)

	if err != nil {
		log.Println(err)
		return ret
	}

	return ret
}

func cliDesc(id string) {
	movie := desc(id)

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintf(writer, "%v\t%s (%v)\n", movie.id, movie.title, movie.releaseDateTime.Year())
	fmt.Fprintf(writer, "Tagline: \t%v\n", movie.tagline)
	fmt.Fprintf(writer, "Original title: \t%v\n", movie.originalTitle)
	fmt.Fprintf(writer, "Runtime: \t%v\n", movie.runtime)
	fmt.Fprintf(writer, "Languages: \t%v\n", movie.spokenLanguages)

	overviewWords := strings.Split(movie.overview, " ")

	fmt.Fprintf(writer, "Overview: \t")

	{
		const maxLen = 50

		var count int
		for _, word := range overviewWords {
			count += len(word)
			if count >= maxLen {
				count = 0
				fmt.Fprintf(writer, "\n \t")
			}

			fmt.Fprintf(writer, "%v ", word)
		}

		fmt.Fprintf(writer, "\n")
	}

	fmt.Fprintf(writer, "Adult rated: \t%v\n", movie.adult)
	fmt.Fprintf(writer, "Budget: \t$%v\n", movie.budget)
	fmt.Fprintf(writer, "Genres: \t%v\n", movie.genres)
	fmt.Fprintf(writer, "Homepage: \t%v\n", movie.homepage)
	fmt.Fprintf(writer, "imdb id: \t%v\n", movie.imdbID)
	fmt.Fprintf(writer, "Languages: \t%v\n", movie.originalLanguage)
	fmt.Fprintf(writer, "Popularity: \t%v\n", movie.popularity)
	fmt.Fprintf(writer, "Production Companies: \t%v\n", movie.productionCompanies)
	fmt.Fprintf(writer, "Production Countries: \t%v\n", movie.productionCountries)
	fmt.Fprintf(writer, "Status: \t%v\n", movie.status)

	fmt.Fprintf(writer, "Video: \t%v\n", movie.video)
	fmt.Fprintf(writer, "Vote Average: \t%v\n", movie.voteAverage)
	fmt.Fprintf(writer, "Vote Count: \t%v\n", movie.voteCount)

	// unsure if useful for now
	// fmt.Fprintf(writer, "", movie.posterPath)
	// fmt.Fprintf(writer, "", movie.revenue)
	// fmt.Fprintf(writer, "Collection: \t%v\n", movie.belongsToCollection)

	writer.Flush()
}
