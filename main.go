package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"gopkg.in/redis.v3"
)

// RedisChatHistoryKey stores the Redis key to the sorted set
// with chat room history
const RedisChatHistoryKey string = "gochat:room:history:"

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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// get the room going
	go r.run(redisClient)

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
