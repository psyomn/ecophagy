package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func acquireKey() byte {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter 1 word: ")
	text, _ := reader.ReadString('\n')

	var ret byte
	for _, c := range text {
		ret += byte(c)
	}
	return ret
}

func applyKey(data []byte, key byte) []byte {
	ret := make([]byte, len(data))
	for i := range data {
		ret[i] = data[i] + key
	}
	return ret
}

func unapplyKey(data []byte, key byte) []byte {
	ret := make([]byte, len(data))
	for i := range data {
		ret[i] = data[i] - key
	}
	return ret
}

func binaryToText(path string, key byte) (string, error) {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	decBin := applyKey(bin, key)

	ret := ""

	for _, byte := range decBin {
		ret += fmt.Sprintf("%x ", byte)
	}

	ret = strings.TrimSuffix(ret, " ")

	return ret, nil
}

func textToBinary(path string) ([]byte, error) {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	str := string(bin)
	values := strings.Split(str, " ")

	original := make([]byte, len(values))
	for i := range values {
		decodedByte, err := strconv.ParseInt(values[i], 16, 64)
		panicIf(err)
		original[i] = byte(decodedByte)
	}

	return original, nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// newExt should be ".blah"
func replaceExt(path, newExt string) string {
	oldExt := filepath.Ext(path)
	newName := strings.TrimSuffix(path, oldExt)
	newName += newExt
	return newName
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	key := acquireKey()
	fmt.Println("your key is:", key)
	fmt.Println("checking in path:", dir)

	var textFiles []string
	var binaryFiles []string

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".txt" {
			textFiles = append(textFiles, path)
		} else {
			binaryFiles = append(binaryFiles, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("text files: ")
	for _, file := range textFiles {
		fmt.Println("-", file)

		bin, err := textToBinary(file)
		panicIf(err)

		encBin := unapplyKey(bin, key)

		err = ioutil.WriteFile(replaceExt(file, ".bin"), encBin, 0600)
		panicIf(err)
	}

	fmt.Println("binary files: ")
	for _, file := range binaryFiles {
		fmt.Println("-", file)

		str, err := binaryToText(file, key)
		panicIf(err)

		err = ioutil.WriteFile(replaceExt(file, ".txt"), []byte(str), 0600)
		panicIf(err)
	}
}
