package models

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

const (
	sqlNoteFindAllByUserID = `
SELECT
	id,title,comment,contents,view_mode,created_at,updated_at,owner_id
FROM
	notes
WHERE
	owner_id = ?`

	sqlNoteFindByIDAndOwnerID = `
SELECT
	id,title,comment,contents,view_mode,created_at,updated_at,owner_id
FROM
	notes
WHERE
	id = ? AND owner_id = ?`

	sqlNoteUpdateByIDAndOwnerID = `
UPDATE notes
SET
	title = ?,
	comment = ?,
	contents = ?,
	updated_at = ?,
	view_mode = ?
WHERE
	id = ? AND owner_id = ?`

	sqlNoteInsert = `
INSERT INTO notes
(
	title,
	comment,
	contents,
	owner_id,
	created_at,
	updated_at,
	view_mode
)
VALUES
	(?,?,?,?,?,?,?)`

	sqlNotesDeleteByIDAndOwnerID = `
DELETE FROM notes
WHERE
	id = ? AND owner_id = ?`

	sqlNotesUpdateViewModeByOwnerIDAndID = `
UPDATE notes
SET
	view_mode = ?
WHERE
	owner_id = ? AND id = ?`

	sqlNotesSelectPublic = `
SELECT
	notes.id,
	notes.title,
	notes.comment,
	notes.created_at,
	notes.updated_at,
	users.name
FROM notes LEFT JOIN users
ON notes.owner_id = users.id
WHERE
	view_mode = 1`

	sqlSelectPublicNote = `
SELECT
	notes.id,
	notes.title,
	notes.comment,
	notes.contents,
	notes.view_mode,
	notes.created_at,
	notes.updated_at,
	notes.owner_id,
	users.name
FROM notes LEFT JOIN users
ON notes.owner_id = users.id
WHERE
	notes.id = ? AND view_mode = 1`
)

type ViewMode int

const (
	ViewModePrivate ViewMode = 0
	ViewModePublic  ViewMode = 1
)

func NotesFindByIDandOwnerID(id, ownerid int) *Note {
	row, err := Handler.QueryRow(sqlNoteFindByIDAndOwnerID, id, ownerid)
	if err != nil {
		log.Println(err)
		return nil
	}
	n := &Note{}

	if err := row.Scan(&n.ID, &n.Title, &n.Comment, &n.Contents,
		&n.ViewMode, &n.CreatedAt, &n.UpdatedAt, &n.OwnerID); err != nil {
		log.Println(err)
		return nil
	}

	return n
}

func NotesPublicFindByID(id int) (*Note, error) {
	row, err := Handler.QueryRow(sqlSelectPublicNote, id)
	if err != nil {
		return nil, err
	}
	n := &Note{}

	if err := row.Scan(
		&n.ID, &n.Title, &n.Comment, &n.Contents, &n.ViewMode,
		&n.CreatedAt, &n.UpdatedAt, &n.OwnerID, &n.Username,
	); err != nil {
		return nil, err
	}

	return n, nil
}

func NotesFindAllByUserID(ownerID int) ([]*Note, error) {
	ret := make([]*Note, 0, 16)

	rows, err := Handler.Query(sqlNoteFindAllByUserID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n := &Note{}

		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Comment,
			&n.Contents,
			&n.ViewMode,
			&n.CreatedAt,
			&n.UpdatedAt,
			&n.OwnerID,
		)

		if err != nil {
			return nil, err
		}

		ret = append(ret, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func NotesAllPublic() ([]*Note, error) {
	ret := make([]*Note, 0, 16)

	rows, err := Handler.Query(sqlNotesSelectPublic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		n := &Note{}

		err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Comment,
			&n.CreatedAt,
			&n.UpdatedAt,
			&n.Username,
		)

		if err != nil {
			return nil, err
		}

		ret = append(ret, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func NotesUpdateByIDAndOwnerID(
	title, comment, contents string,
	updatedAt time.Time,
	viewMode int,
	id, ownerID int,
) (*sql.Result, error) {
	res, err := Handler.Execute(
		sqlNoteUpdateByIDAndOwnerID,
		title,
		comment,
		contents,
		updatedAt,
		viewMode,
		id,
		ownerID,
	)
	return res, err
}

func NotesUpdateViewModeByOwnerIDAndID(
	viewMode ViewMode,
	ownerID, id int,
) (*sql.Result, error) {
	res, err := Handler.Execute(
		sqlNotesUpdateViewModeByOwnerIDAndID,
		viewMode,
		ownerID,
		id,
	)
	return res, err
}

func NotesInsert(
	title, comment, contents string,
	createdAt, updatedAt time.Time,
	viewMode, ownerID int,
) (*sql.Result, error) {
	res, err := Handler.Execute(
		sqlNoteInsert,
		title,
		comment,
		contents,
		ownerID,
		createdAt,
		updatedAt,
		viewMode,
	)
	return res, err
}

func ViewModeFromStrOrDefault(vs string, d int) (ret int, usedDefault bool) {
	if vi, err := strconv.Atoi(vs); err != nil {
		ret = d
		usedDefault = true
	} else {
		ret = vi
		usedDefault = false
	}

	return ret, usedDefault
}

func NotesDeleteByIDAndOwnerID(id, ownerID int) (*sql.Result, error) {
	res, err := Handler.Execute(
		sqlNotesDeleteByIDAndOwnerID,
		id, ownerID,
	)
	return res, err
}

type Note struct {
	ID        uint
	OwnerID   int
	Title     string
	Comment   string
	Contents  string
	ViewMode  int
	CreatedAt time.Time
	UpdatedAt time.Time

	// auxiliary fields
	Username *string
}
