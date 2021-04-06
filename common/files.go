package common

import (
	// TODO: go 1.16 should just use io (see go doc ioutil)
	"io/ioutil"
	"os"
	"path"
)

func FileToBytes(filename string) ([]byte, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fs.Close()

	bytes, err := ioutil.ReadAll(fs)
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
