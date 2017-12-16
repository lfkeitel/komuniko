package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	rooms map[string]*room
}

func New() *Server {
	return &Server{
		rooms: make(map[string]*room),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	urlparts := strings.SplitN(r.URL.Path, "/", 4)
	if len(urlparts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	roomName := urlparts[3]
	chatroom, exists := s.rooms[roomName]
	if !exists {
		chatroom = newRoom(roomName)
		s.rooms[roomName] = chatroom
		go chatroom.relay()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	chatter := &client{
		room: chatroom,
		conn: conn,
		send: make(chan *message, 256),
	}

	chatroom.register(chatter)

	go chatter.writePump()
	go chatter.readPump()
}
