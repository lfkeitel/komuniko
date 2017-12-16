package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/lfkeitel/komuniko/server"
)

var (
	configFilePath string
	addr           string
)

func init() {
	flag.StringVar(&configFilePath, "c", "config.toml", "Configuration file")
	flag.StringVar(&addr, "addr", ":8080", "HTTP address")
}

func main() {
	flag.Parse()

	server := server.New()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/room/", serverRoom)
	http.HandleFunc("/ws/room/", server.ServeWs)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.Redirect(w, r, "/room/main", http.StatusTemporaryRedirect)
}

func serverRoom(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}
