package server

import (
	"fmt"
	"time"
)

type room struct {
	clients   map[*client]bool
	name      string
	broadcast chan *message
}

func newRoom(name string) *room {
	return &room{
		name:      name,
		clients:   make(map[*client]bool),
		broadcast: make(chan *message),
	}
}

func (r *room) register(c *client) {
	r.clients[c] = true
	c.send <- newMessage("Server", fmt.Sprintf("Welcome to chatroom %s", r.name))
	r.broadcast <- newMessage("Server", fmt.Sprintf("%s has joined", c.conn.RemoteAddr().String()))
}

func (r *room) unregister(c *client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		close(c.send)
		r.broadcast <- newMessage("Server", fmt.Sprintf("%s has left", c.conn.RemoteAddr().String()))
	}
}

func (r *room) relay() {
	for {
		select {
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

type message struct {
	User      string `json:"user"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func newMessage(user, msg string) *message {
	return &message{
		User:      user,
		Message:   msg,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
