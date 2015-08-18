package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"strings"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Contains is an InArray working with strings
func Contains(list []string, elem string) bool {
	for _, t := range list {
		if strings.Contains(elem, t) {
			return true
		}
	}
	return false
}

// RandomYoutubeVideo returns a random YouTube video
func RandomYoutubeVideo() string {
	return "http://www.youtube.com"
}
