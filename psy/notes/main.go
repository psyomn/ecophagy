package notes

import (
	"flag"
	"fmt"
	"log"
	"path"
	"strings"

	psycommon "git.sr.ht/~psyomn/ecophagy/psy/common"
	"git.sr.ht/~psyomn/ecophagy/psy/notes/models"
	"git.sr.ht/~psyomn/ecophagy/psy/notes/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Register struct {
	UsernamePassword string // USERNAME:PASSWORD
	Username         string
	Password         string
}

func (s *Register) Process() error {
	parts := strings.Split(s.UsernamePassword, ":")

	if len(parts) < 2 {
		return ErrBadRegisterCLI
	}

	s.Username = parts[0]
	s.Password = parts[1]

	const (
		usz = 3
		psz = 10
	)

	if len(s.Username) < usz {
		log.Println("username must be at least", usz, "characters")
		return fmt.Errorf("%w: username", ErrBadSize)
	}

	if len(s.Password) < psz {
		log.Println("password must be at least", psz, "characters")
		return fmt.Errorf("%w: password", ErrBadSize)
	}

	_, err := models.UserCreate(s.Username, s.Password)
	return err
}

type Session struct {
	Server   Server
	Register Register
}

func sessionFromArgs(sess *Session, args []string) *Session {
	fs := flag.NewFlagSet("notes", flag.ExitOnError)

	fs.StringVar(&sess.Server.Host, "host", "0.0.0.0", "set the host")
	fs.StringVar(&sess.Server.Port, "port", "15000", "set the port")
	fs.StringVar(&sess.Server.DataDirPath, "data", ".", "set the data path/directory")
	fs.StringVar(&sess.Register.UsernamePassword, "register", "", "register a user via USERNAME:PASSWORD format")

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	return sess
}

func Run(args psycommon.RunParams) psycommon.RunReturn {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmsgprefix)

	log.Println("creating database")

	sess := sessionFromArgs(&Session{
		Server: Server{Cache: make(map[SessionID]*SessionData)},
	}, args)

	log.Println("creating database connection...")

	h, err := models.HandleNew(path.Join(sess.Server.DataDirPath, storage.Name))
	if err != nil {
		panic(fmt.Sprintf("could not create the database: %v", err))
	}
	models.Handler = h

	if err := storage.MaybeCreateDB(h.GetRaw()); err != nil {
		log.Println("problem creating the database:", storage.Name, "error:", err)
		return err
	}

	if sess.Register.UsernamePassword != "" {
		log.Println("registering user...")
		return sess.Register.Process()
	}

	log.Println("starting with session vars:", sess)

	return NewServer(sess).Srv.ListenAndServe()
}
