package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type movieRecord struct {
	adult               bool
	belongsToCollection string
	budget              uint64
	genres              string
	homepage            string
	id                  uint64
	imdbID              uint64
	originalLanguage    string
	originalTitle       string
	overview            string
	popularity          float64
	posterPath          string
	productionCompanies string
	productionCountries string
	releaseDate         int64
	releaseDateTime     time.Time // for db driver
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

	// This is python dict in the data set:
	//   [{'id': 28, 'name': 'Action'},
	//    {'id': 18, 'name': 'Drama'},
	//    {'id': 53, 'name': 'Thriller'}]
	s.genres = minifyNamePyDict(line[3])
	s.homepage = line[4]
	{
		id, err := strconv.ParseUint(line[5], 10, 64)
		if err == nil {
			s.id = id
		}
	}

	s.imdbID = parseImdbID(line[6])
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

	//   [{'name': 'Miramax Films', 'id': 14},
	//    {'name': 'A Band Apart', 'id': 59}]
	s.productionCompanies = minifyNamePyDict(line[12])

	//   [{'iso_3166_1': 'GB',
	//     'name': 'United Kingdom'}]
	s.productionCountries = minifyNamePyDict(line[13])
	s.releaseDate = parseMovieDate(line[14])

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

	//  [{'iso_639_1': 'en',
	//    'name': 'English'}]
	s.spokenLanguages = minifyNamePyDict(line[17])

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

func parseMovieDate(date string) int64 {
	const layout = "2006-01-02"

	if len(date) < len(layout) {
		return 0
	}

	t, err := time.Parse(layout, date)
	if err != nil {
		log.Println(err, "date:", date)
		return 0
	}
	return t.Unix()
}

func parseImdbID(id string) uint64 {
	// ids are of the form tt123123, so we just get rid of the t's
	// and store ints for more compact space...

	if len(id) < 2 {
		log.Println("movie has no imdb id")
		return 0
	}

	convID, err := strconv.ParseUint(id[2:], 10, 64)
	if err != nil {
		log.Println(err, "id:", id)
		return 0
	}
	return convID
}

// TODO: handle utf8 stuff properly. See output in tests for examples.
func minifyNamePyDict(line string) string {
	var cursor int
	var names []string
	for {
		name, next := parseEntry(line[cursor:])
		cursor += next

		if name == "" {
			break
		}

		names = append(names, name)
	}

	return strings.Join(names, ",")
}

// this is necessary because the data has troublesome, python dictionary
// entries. Take a look at the test cases for further elaboration
func parseEntry(line string) (string, int) {
	lookup := `'name': `
	startIndex := strings.Index(line, lookup)

	if startIndex < 0 {
		// not found, error
		return "", startIndex
	}

	offset := startIndex + len(lookup)

	cursor := offset
	terminator := line[cursor]

	cursor++ // move to start of input

	var retName string

	for {
		terminated := line[cursor] == terminator
		isWithinBounds := cursor+1 < len(line)-1
		shouldStop := terminated && isWithinBounds &&
			(line[cursor+1] == ',' || line[cursor+1] == '}')

		if shouldStop {
			break
		}

		retName += string(line[cursor])
		cursor++
	}

	if len(retName) >= 2 &&
		retName[0] == terminator &&
		retName[len(retName)-1] == terminator {
		// Some python dictionary entries can have docsrings instead
		// so we can do this rather dodgy check, and get rid of the
		// surplus terminators, without making the parser more
		// complicated...
		retName = retName[1 : len(retName)-1]
	}

	return retName, cursor
}
