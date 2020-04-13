package main

const (
	createMovieMetadataTable = `
CREATE TABLE IF NOT EXISTS movies (
  adult                   BOOLEAN,
  belongs_to_collection   STRING, -- python obj
  budget                  BIGINT,
  genres                  TEXT,
  homepage                TEXT,
  id                      INTEGER PRIMARY KEY,
  imdb_id                 BIGINT,
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

	createUserRatingsTable = `
CREATE TABLE IF NOT EXISTS ratings (
  id      INTEGER PRIMARY KEY AUTOINCREMENT,
  comment TEXT,
  score   INTEGER -- should be within range of [1, 10]
);
`
	insertMovie = `INSERT INTO movies values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	createMovieIDIndex = `CREATE INDEX IF NOT EXISTS movie_id_index ON movies(id);`

	createMovieNameIndex = `CREATE INDEX IF NOT EXISTS movie_original_title_index ON movies(original_title)`

	// NB: prepared statement must also lowercase the input
	searchMovieIlike = `SELECT * FROM movies WHERE LOWER(title) like '%?%'`
)
