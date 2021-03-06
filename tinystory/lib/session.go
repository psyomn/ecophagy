package tinystory

type Session struct {
	Host               string
	Port               string
	Repository         string
	Assets             string
	ExperimentalParser bool
}

func MakeDefaultSession() *Session {
	return &Session{
		Host:       "127.0.0.1",
		Port:       "9090",
		Repository: "./stories",
		Assets:     "./assets",
	}
}
