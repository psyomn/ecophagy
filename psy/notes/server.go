package notes

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/psyomn/ecophagy/psy/notes/models"
	"github.com/psyomn/ecophagy/psy/notes/static"
)

const CookieSessionName = "Psy-Notes-Session"

var staticPages = static.PagesNew()

func CookieFind(cookies []*http.Cookie, name string) *http.Cookie {
	ckix := slices.IndexFunc(cookies, func(c *http.Cookie) bool {
		return name == c.Name
	})

	if ckix >= 0 {
		return cookies[ckix]
	}

	return nil
}

func GenerateToken() ([]byte, error) {
	bs := make([]byte, 128)
	_, err := rand.Read(bs)
	if err != nil {
		return nil, errors.Join(ErrGenerateToken, err)
	}
	return bs, nil
}

type SessionID string

type SessionData struct {
	UserID   int
	Username string
}

type Server struct {
	Host        string
	Port        string
	DataDirPath string
	CacheMx     sync.RWMutex
	Cache       map[SessionID]*SessionData
	Srv         *http.Server
}

type PageUser struct {
	Name string
}

type PageView struct {
	User         *PageUser
	Notes        []*models.Note
	IsPublicNote bool
}

func RedirectOnErr(w http.ResponseWriter, r *http.Request, err error, maybePath ...string) error {
	if err == nil {
		return nil
	}

	log.Println("error:", err)

	path := "/"
	if len(maybePath) > 0 {
		path = maybePath[0]
	}

	http.Redirect(w, r, path, http.StatusSeeOther)

	return err
}

func ExpireCookie(c *http.Cookie, w http.ResponseWriter) {
	c.Value = ""
	c.Path = "/"
	c.MaxAge = -1
	c.HttpOnly = true
	http.SetCookie(w, c)
}

func (s *Server) SessionDataFromCookies(w http.ResponseWriter, r *http.Request) *SessionData {
	cks := r.Cookies()

	ix := slices.IndexFunc(cks, func(c *http.Cookie) bool { return c.Name == CookieSessionName })
	if ix == -1 {
		return nil
	}

	ck := cks[ix]

	s.CacheMx.RLock()
	defer s.CacheMx.RUnlock()
	v, ok := s.Cache[SessionID(ck.Value)]
	if !ok {
		/* delete cookie if we don't know about them */
		ExpireCookie(ck, w)
		return nil
	}

	return v
}

func (s *Server) SessionClear(w http.ResponseWriter, r *http.Request) {
	cs := CookieFind(r.Cookies(), CookieSessionName)
	if cs == nil {
		return
	}

	s.CacheMx.Lock()
	defer s.CacheMx.Unlock()

	delete(s.Cache, SessionID(cs.Value))

	/* blah!  raisins! */
	ExpireCookie(cs, w)
}

func (s *Server) HandleDefault(w http.ResponseWriter, r *http.Request) {
	pd := PageView{User: nil, Notes: []*models.Note{}}

	sd := s.SessionDataFromCookies(w, r)
	if sd != nil {
		pd.User = &PageUser{Name: sd.Username}
	}

	n, err := models.NotesAllPublic()
	if RedirectOnErr(w, r, err) != nil {
		return
	}

	pd.Notes = n

	w.WriteHeader(http.StatusOK)
	if err := staticPages.Index.Execute(w, &pd); err != nil {
		log.Println("error serving page:", err)
	}
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	s.SessionClear(w, r)

	if err := r.ParseForm(); err != nil {
		log.Println("form-parse:", err)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	usr, err := models.UserLogin(username, password)
	if RedirectOnErr(w, r, err) != nil {
		return
	}

	bs, err := GenerateToken()
	if RedirectOnErr(w, r, err) != nil {
		return
	}

	s.CacheMx.Lock()
	defer s.CacheMx.Unlock()
	sid := fmt.Sprintf("%x", bs)
	s.Cache[SessionID(sid)] = &SessionData{
		UserID:   usr.ID,
		Username: usr.Name,
	}

	http.SetCookie(w, &http.Cookie{
		Name:   CookieSessionName,
		Value:  sid,
		MaxAge: 48 * 60 * 60,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	s.SessionClear(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) HandleNotesCreate(
	w http.ResponseWriter,
	r *http.Request,
	pd *PageView,
	sd *SessionData,
	action string,
) {
	switch action {
	case "new":
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			pd.Notes = []*models.Note{{}}
			if err := staticPages.NotesEdit.Execute(w, &pd); err != nil {
				log.Println("error serving page:", err)
				return
			}
		case "POST":
			if err := r.ParseForm(); err != nil {
				log.Println("form-parse:", err)
				return
			}

			title := r.PostForm.Get("title")
			comment := r.PostForm.Get("comment")
			contents := r.PostForm.Get("contents")

			viewModeRaw := r.PostForm.Get("view_mode")
			viewMode, ud := models.ViewModeFromStrOrDefault(viewModeRaw, 0)
			if ud {
				log.Println("warn: used default for view mode because view_mode=", viewModeRaw)
			}

			if _, err := models.NotesInsert(
				title, comment, contents, time.Now(), time.Now(),
				viewMode, sd.UserID); err != nil {
				log.Println("error inserting notes:", err)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "/notes", http.StatusSeeOther)
			return

		default:
			log.Println("new: unsupported method:", r.Method)
		}
	default:
		w.WriteHeader(http.StatusOK)
		if err := staticPages.Notes.Execute(w, &pd); err != nil {
			log.Println("error serving page:", err)
		}
	}
}

func (s *Server) HandleNotesRUD(
	w http.ResponseWriter,
	r *http.Request,
	pd *PageView,
	sd *SessionData,
	action string,
	id int,
) {

	switch action {
	case "view":
		n := models.NotesFindByIDandOwnerID(id, sd.UserID)
		pd.Notes = []*models.Note{n}

		w.WriteHeader(http.StatusOK)
		if err := staticPages.NotesView.Execute(w, pd); err != nil {
			log.Println(err)
		}
		return

	case "edit":
		n := models.NotesFindByIDandOwnerID(id, sd.UserID)
		pd.Notes = []*models.Note{n}

		w.WriteHeader(http.StatusOK)
		if err := staticPages.NotesEdit.Execute(w, pd); err != nil {
			log.Println(err)
		}
		return

	case "destroy":
		if _, err := models.NotesDeleteByIDAndOwnerID(id, sd.UserID); err != nil {
			log.Println("error deleting note:", err)
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return

	case "update":
		if err := RedirectOnErr(w, r, r.ParseForm()); err != nil {
			return
		}

		viewMode, _ := models.ViewModeFromStrOrDefault(r.FormValue("view_mode"), 0)

		_, err := models.NotesUpdateByIDAndOwnerID(
			r.FormValue("title"),
			r.FormValue("comment"),
			r.FormValue("contents"),
			time.Now(),
			viewMode,
			id,
			sd.UserID,
		)

		viewPath := fmt.Sprintf("/notes/view/%d", id)
		if RedirectOnErr(w, r, err, viewPath) != nil {
			log.Println("could not update row:", err)
			return
		}

		http.Redirect(w, r, viewPath, http.StatusSeeOther)
		return
	case "publish":
		if _, err := models.NotesUpdateViewModeByOwnerIDAndID(
			models.ViewModePublic, sd.UserID, id,
		); RedirectOnErr(w, r, err, "/notes") != nil {
			return
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	case "hide":
		if _, err := models.NotesUpdateViewModeByOwnerIDAndID(
			models.ViewModePrivate, sd.UserID, id,
		); RedirectOnErr(w, r, err, "/notes") != nil {
			return
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	default:
		log.Println("unknown action", action)
	}
}

/**
 * GET  /notes/
 * GET  /notes/new -- render view
 * POST /notes/new -- consume data for new note
 * GET  /notes/view/1
 * GET  /notes/edit/1
 * POST /notes/update/1 -- forms don't support anything apart from GET/POST
 * GET  /notes/destroy/1
 * GET  /notes/publish/1
 * GET  /notes/hide/1
 */
func (s *Server) HandleNotes(w http.ResponseWriter, r *http.Request) {
	checkPathFn := func(p string) ([]string, error) {
		parts := strings.Split(p, "/")
		if len(parts) < 3 {
			return nil, ErrNotesBadPath
		}
		return parts, nil
	}

	pd := PageView{User: nil, Notes: []*models.Note{}}
	sd := s.SessionDataFromCookies(w, r)
	if sd == nil {
		/* must be logged in */
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pd.User = &PageUser{Name: sd.Username}
	userNotes, err := models.NotesFindAllByUserID(sd.UserID)
	if RedirectOnErr(w, r, err) != nil {
		return
	}
	pd.Notes = userNotes

	parts, err := checkPathFn(r.URL.Path)
	if RedirectOnErr(w, r, err) != nil {
		return
	}
	action := parts[2]

	switch action {
	case "new":
		s.HandleNotesCreate(w, r, &pd, sd, action)
		return
	case "view", "edit", "update", "destroy", "publish", "hide":
		id, err := strconv.Atoi(parts[3])
		if RedirectOnErr(w, r, err) != nil {
			return
		}

		s.HandleNotesRUD(w, r, &pd, sd, action, id)
		return
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

/* GET /public/1 */
func (s *Server) HandlePublic(w http.ResponseWriter, r *http.Request) {
	checkPathFn := func(p string) ([]string, error) {
		parts := strings.Split(p, "/")
		if len(parts) < 3 {
			return nil, fmt.Errorf("%w: bad size", ErrPublicBadPath)
		}
		return parts, nil
	}

	pd := PageView{User: nil, Notes: []*models.Note{}}
	sd := s.SessionDataFromCookies(w, r)
	if sd != nil {
		pd.User = &PageUser{Name: sd.Username}
	}

	parts, err := checkPathFn(r.URL.Path)
	if RedirectOnErr(w, r, err) != nil {
		return
	}

	id, err := strconv.Atoi(parts[2])
	if RedirectOnErr(w, r, err) != nil {
		return
	}

	n, err := models.NotesPublicFindByID(id)
	if RedirectOnErr(w, r, err) != nil {
		return
	}

	pd.Notes = []*models.Note{n}
	w.WriteHeader(http.StatusOK)
	if err := staticPages.NotesView.Execute(w, pd); err != nil {
		log.Println("error rendering public note:", err)
		return
	}
}

func NewServer(sess *Session) *Server {
	mux := http.NewServeMux()

	S := &Server{
		Cache:   make(map[SessionID]*SessionData),
		CacheMx: sync.RWMutex{},
		Srv: &http.Server{
			Addr:              net.JoinHostPort(sess.Server.Host, sess.Server.Port),
			ReadHeaderTimeout: time.Second * 30,
			Handler:           mux,
		},
	}

	mux.HandleFunc("/public/", S.HandlePublic)
	mux.HandleFunc("/notes/", S.HandleNotes)
	mux.HandleFunc("/login", S.HandleLogin)
	mux.HandleFunc("/logout", S.HandleLogout)
	mux.HandleFunc("/", S.HandleDefault)

	return S
}
