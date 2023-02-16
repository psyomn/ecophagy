package common

import (
	"io"
	"os"
	"path"
)

func FileToBytes(filename string) ([]byte, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	bytes, err := io.ReadAll(fs)
	return bytes, err
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err != nil
}

// ConfigDir prefers $HOME/.config, regardles of XDG stuff (for now)
func ConfigDir() string {
	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		panic("need home to run")
	}

	return path.Join(homeDir, ".config", "psy")
}
