package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	hub := NewHub()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", WebSocketHandler(hub))

	fmt.Println("🌐 Web UI → http://localhost:8080")
	fmt.Println("🔌 Listener WS → ws://localhost:8080/ws")
	fmt.Println("🔌 Control WS → ws://localhost:8080/ws?password=")

	log.Fatal(http.ListenAndServe(Port, nil))
}
