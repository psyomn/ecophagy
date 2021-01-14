package common

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const orgName = "ecophagy"

// FileExists checks if a file exists
func FileExists(filename string) bool {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer fd.Close()
	return err != nil
}

// DownloadIfNotExist will check for a filename, and if it doesn't
// exist, will attempt to download it by a provided url.
func DownloadIfNotExist(filename, url string) error {
	if FileExists(filename) {
		return nil
	}

	return DownloadFile(filename, url)
}

// DownloadFile will download a file from a url
func DownloadFile(filename, url string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// FileToLines will fully read a file, and return it as individual
// lines.
func FileToLines(filename string) ([]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	var line []byte
	var ret []string

	for {
		line, _, err = reader.ReadLine()
		if err != nil {
			break
		}

		ret = append(ret, string(line))
	}

	return ret, nil
}

// HasHome will check if the HOME environment variable is set
func HasHome() bool { return os.Getenv("HOME") != "" }

// HasEditor will check if there is an editor set
func HasEditor() bool { return os.Getenv("EDITOR") != "" }

// Editor will return the editor environment variable, "EDITOR"
func Editor() string { return os.Getenv("EDITOR") }

// DataPath will return the path of the ecophagy project config
// folder. Subprojects may be stored in the form of:
// - $HOME/.config/ecophagy/mvpic
// - $HOME/.config/ecophay/randparty
func DataPath() string {
	return path.Join(os.Getenv("HOME"), ".config", orgName)
}

// PathExists will check if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

// FileToMd5Sum will calculate the hash of a particular file
// contents. The intended use is for quick checks rather than
// secure.
func FileToMd5Sum(path string) (string, error) {
	fh, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, fh); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// FileList will return a list of files, by walking through a
// directory tree.
func FileList(dirpath string) ([]string, error) {
	var ret []string

	err := filepath.Walk(dirpath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			ret = append(ret, path)
			return nil
		})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Newline returns the newline symbol accepted as a newline symbol
func Newline() string { return "\n" }

func GetUserName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "default", err
	}

	return currentUser.Username, nil
}
