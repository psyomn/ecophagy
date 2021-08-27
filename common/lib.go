package common

import (
	"bufio"
	"context"
	"time"

	// nolint:gosec // see below usage
	"crypto/md5"

	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const (
	Newline = '\n'
	orgName = "ecophagy"
)

// DownloadIfNotExist will check for a filename, and if it doesn't
// exist, will attempt to download it by a provided url.
func DownloadIfNotExist(filename, url string) error {
	if FileExists(filename) {
		return nil
	}

	return DownloadFile(filename, url)
}

// DownloadFile will download a file from a url to a designated path
func DownloadFile(path, url string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
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

	for ; !errors.Is(err, io.EOF); line, _, err = reader.ReadLine() {
		if errors.Is(err, io.EOF) {
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
	return !os.IsNotExist(err)
}

// FileToMd5Sum will calculate the hash of a particular file
// contents. The intended use is for quick checks rather than
// secure. Might need to microbenchmark sha256 and verify.
func FileToMd5Sum(path string) (string, error) {
	fh, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	// nolint:gosec // security not important; speed is
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

func GetUserName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "default", err
	}

	return currentUser.Username, nil
}
