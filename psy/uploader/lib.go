/*
Package uploader has all the logic required to spin up an upload http
server, to send files from one computer to another. This is designed
for ease of use with family members, and should only be used in home
networks.

A good portion of the upload code was taken, and repurposed from here:
https://astaxie.gitbooks.io/build-web-application-with-golang/en/04.5.html
*/
package uploader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/psyomn/ecophagy/psy/common"
)

const (
	uploadFileHTML = `<html>
<head>
       <title>Upload file</title>
</head>
<body>
<h1> Your possible (local network) IPs </h1>
<p> {{.IPStr}} </p>
<form enctype="multipart/form-data" action="upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="submit" value="upload" />
</form>
</body>
</html>
`
	uploadsDir = "uploads/"

	port = ":9090"
)

// Run will run the command with default configs. For now, the
// uploader does not accept any configuration.
func Run(_ common.RunParams) common.RunReturn {
	createDirs()

	fmt.Println("listening at port:", port)

	ips, _ := common.GetLocalIP()
	fmt.Println("your possible IPs: ")
	for _, ip := range ips {
		fmt.Println(" ", ip.To4().String())
	}

	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(port, nil))

	return nil
}

func createDirs() {
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			log.Println("could not create uploads dir: ", err)
			os.Exit(1)
		}
		log.Println("created uploads dir")
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	type homepage struct {
		IPStr string
	}

	if r.Method == "GET" {
		var ips []string
		ipObjs, _ := common.GetLocalIP()
		for _, ip := range ipObjs {
			ips = append(ips, ip.To4().String())
		}
		ipStr := strings.Join(ips, ",")
		uploadT := template.Must(template.New("upload-page").Parse(uploadFileHTML))

		var buff bytes.Buffer
		buffw := bufio.NewWriter(&buff)
		if err := uploadT.Execute(buffw, &homepage{IPStr: ipStr}); err != nil {
			log.Println("problem rendering template:", err)
		}
		buffw.Flush() // for some reason, need to flush explicitly

		if _, err := w.Write(buff.Bytes()); err != nil {
			log.Println("problem writing webpage", err)
			return
		}

		return
	}

	if r.Method != "POST" {
		_, err := w.Write([]byte("only supports POST and GET"))
		if err != nil {
			log.Println("error writing response", err)
		}
		return
	}

	if err := r.ParseMultipartForm(1 << 27); err != nil {
		log.Println("problem parsing multipart form", err)
		return
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	uploadsPath := filepath.Join(uploadsDir, handler.Filename)

	f, err := os.OpenFile(
		uploadsPath,
		os.O_WRONLY|os.O_CREATE,
		0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		log.Println("problem uploading file:", err)
		return
	}

	log.Println("uploaded file: ", handler.Filename)
}
