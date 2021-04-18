package tinystory

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/psyomn/ecophagy/common"
)

type Server struct {
	visitor    *Visitor
	httpServer *http.Server

	indexTemplate *template.Template
	storyTemplate *template.Template
}

const indexFilename = "index.html"
const storyFilename = "story.html"

func ServerNew(sess *Session, documents []Document) (*Server, error) {
	muxer := http.NewServeMux()

	var indexTemplate *template.Template
	{
		data, err := common.FileToBytes(path.Join(sess.Assets, indexFilename))
		if err != nil {
			return nil, err
		}

		maybeIndex, err := template.New("index").Parse(string(data))
		if err != nil {
			return nil, err
		}
		indexTemplate = maybeIndex
	}

	var storyTemplate *template.Template
	{
		data, err := common.FileToBytes(path.Join(sess.Assets, storyFilename))
		if err != nil {
			return nil, err
		}

		maybeStory, err := template.New("story").Parse(string(data))
		if err != nil {
			return nil, err
		}
		storyTemplate = maybeStory
	}

	server := &Server{
		indexTemplate: indexTemplate,
		storyTemplate: storyTemplate,
		httpServer: &http.Server{
			Addr:    sess.Host + ":" + sess.Port,
			Handler: muxer,
		},
		visitor: VisitorNew(documents),
	}

	muxer.HandleFunc("/", server.HandleRoot)
	muxer.HandleFunc("/story/", server.HandleStory)

	return server, nil
}

func (s *Server) Start() error {
	log.Println("starting server...")
	return s.httpServer.ListenAndServe()
}

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	listing := struct {
		Items []IndexListing
	}{
		Items: s.visitor.GetIndexListing(),
	}

	if err := s.indexTemplate.Execute(w, listing); err != nil {
		fmt.Println("error writing template: ", err)
	}
}

func (s *Server) HandleStory(w http.ResponseWriter, r *http.Request) {
	var storyIndex, nodeIndex int
	{
		thinPath := strings.TrimPrefix(r.URL.Path, "/story/")
		parts := strings.Split(thinPath, "/")

		if len(parts) < 2 {
			renderError(w, "badly formed path")
			return
		}

		maybeStoryIndex, err := strconv.Atoi(parts[0])
		if err != nil {
			renderError(w, "stories should be numeric")
			return
		}

		maybeNodeIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			renderError(w, "parts should be numeric")
			return
		}

		storyIndex = maybeStoryIndex
		nodeIndex = maybeNodeIndex
	}

	responseData := struct {
		Title      string
		Authors    []string
		Website    string
		Fragment   StoryFragment
		StoryIndex int
		NodeIndex  int
	}{
		Title:      s.visitor.Documents[storyIndex].Title,
		Authors:    s.visitor.Documents[storyIndex].Authors,
		Website:    s.visitor.Documents[storyIndex].Website,
		Fragment:   s.visitor.Documents[storyIndex].Fragments[nodeIndex],
		StoryIndex: storyIndex,
	}

	if err := s.storyTemplate.Execute(w, responseData); err != nil {
		renderError(w, "problem getting that fragment")
		return
	}
}

func renderError(w http.ResponseWriter, str string) {
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write([]byte(str)); err != nil {
		fmt.Printf("problem writing error: %s\n", err.Error())
	}
}
