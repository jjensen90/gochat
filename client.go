package main

import (
	"github.com/gorilla/websocket"
	"fmt"
)

// client represents a single chatting user.
type client struct {
	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan []byte

	// room is the room this client is chatting in.
	room *room

	name string
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			stringMessage := string(msg[:])
			if Contains(c.room.commands, stringMessage) {
				c.room.DoCommand(stringMessage, c)
			} else {
				name := fmt.Sprintf("%s: ", c.name)
				c.room.forward <- append([]byte(name), msg...)
			}
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
