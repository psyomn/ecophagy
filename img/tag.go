/*
Package img is just a wrapper to use around whatever can extract and
set exit comment tags.

Right now I decided to use exiftool from the system, as most exif
libraries I've checked online for go seem to be a little inconvenient
to use.
*/
package img

import (
	"fmt"
	"os/exec"
)

func HasExifTool() bool {
	_, err := exec.LookPath("exiftool")
	return err == nil
}

// Wrapper for actual command:
//
//	exiftool -s -s -s -comment picture.jpg
func GetExifComment(filename string) (string, error) {
	out, err := exec.Command(
		"exiftool", "-s", "-s", "-s",
		"-comment", filename,
	).Output()

	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Wrapper for actual command:
//
//	exiftool -comment="blargh" picture.jpg
func SetExifComment(filename string, comment string) error {
	cmdComment := fmt.Sprintf("-comment=%s", comment)
	_, err := exec.Command("exiftool", cmdComment, filename).Output()
	return err
}
