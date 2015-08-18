package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"io/ioutil"
	"net/url"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {

	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool

	commands []string
}

// newRoom makes a new room that is ready to go.
func newRoom() *room {
	return &room{
		forward:  make(chan []byte),
		join:     make(chan *client),
		leave:    make(chan *client),
		clients:  make(map[*client]bool),
		commands: []string{"/roll", "/ascii", "/yt", "/yomama"},
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			client.send <- []byte(fmt.Sprintf("Available room commands: %v", r.commands))
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send the message
				default:
					// failed to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
		name:   randomdata.SillyName(),
	}

	r.join <- client

	msg := []byte(fmt.Sprint("A new user has joined: ", client.name))
	r.forward <- msg

	defer func() {
		r.leave <- client
		msg := []byte(fmt.Sprint("A user has left: ", client.name))
		r.forward <- msg
	}()

	go client.write()
	client.read()
}

func (r *room) DoCommand(command string, c *client) bool {
	switch {
	case command == "/roll":
		r.forward <- []byte(fmt.Sprintf("Casino: %s rolled a %d", c.name, rand.Intn(100)))
		return true
	case strings.Contains(command, "/ascii"):
		response, err := r.Asciify(command)
		if err != nil {
			log.Printf("Error asciifying")
		}
		r.forward <- []byte(response)
	case command == "/yt":
		r.forward <- []byte(fmt.Sprintf("YouTube video: %s", GetRandomVideo()))
	case command == "/yomama":
		r.forward <- []byte(fmt.Sprintf("%s", GetYoMamaJoke()))
	default:
		return false
	}
	return true
}

func (r *room) Asciify(command string) (string, error) {
	command = strings.Replace(command, "/ascii", "", 1)
	if command != "" {
		safeString := url.QueryEscape(strings.Trim(command, " "))
		return r.GetAscii(safeString), nil
	}
	return "", fmt.Errorf("Error Asciifying String")
}

func (r *room) GetAscii(query string) string {
	url := "http://artii.herokuapp.com/make?text="
	resp, err := http.Get(fmt.Sprintf("%s%s", url, query))
	if err != nil {
		log.Printf("GetAscii Error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return fmt.Sprintf("<pre>%s</pre>", string(body))
}
