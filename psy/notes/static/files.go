package static

import (
	"embed"
	"fmt"
	"html/template"
	"time"
)

var _ embed.FS

//go:embed head.html
var headRaw string

//go:embed index.html
var indexRaw string

//go:embed notes.html
var notesRaw string

//go:embed notes-actions.html
var notesActionsRaw string

//go:embed notes-view.html
var notesViewRaw string

//go:embed notes-edit.html
var notesEditRaw string

//go:embed menus.html
var menusRaw string

type Pages struct {
	Head  *template.Template
	Index *template.Template
	Menus *template.Template

	Notes     *template.Template
	NotesView *template.Template
	NotesEdit *template.Template
}

func PagesNew() *Pages {
	fnmap := template.FuncMap{
		"hdate": func(t time.Time) string {
			return fmt.Sprintf("%v", t.Format(time.DateTime))
		},
	}

	mustParse := func(s string, t *template.Template) *template.Template {
		tt, err := t.Parse(s)
		if err != nil {
			panic(err)
		}
		return tt
	}

	indexTmpl := template.New("index")
	menusTmpl := template.New("menus")
	headTmpl := template.New("head")

	notesTmpl := template.New("notes")
	notesViewTmpl := template.New("notes-view")
	notesEditTmpl := template.New("notes-edit")

	for _, t := range []*template.Template{
		indexTmpl,
		menusTmpl,
		headTmpl,

		notesTmpl,
		notesViewTmpl,
		notesEditTmpl,
	} {
		t.Funcs(fnmap)
	}

	return &Pages{
		Index: mustParse(
			headRaw+indexRaw+menusRaw,
			indexTmpl),

		Notes: mustParse(
			headRaw+notesRaw+notesActionsRaw+menusRaw,
			notesTmpl),

		NotesView: mustParse(
			headRaw+notesViewRaw+notesActionsRaw+menusRaw,
			notesViewTmpl),

		NotesEdit: mustParse(
			headRaw+notesEditRaw+notesActionsRaw+menusRaw,
			notesEditTmpl),
	}
}
