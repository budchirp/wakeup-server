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

	http.HandleFunc("/ring", RingHandler(hub))
	http.HandleFunc("/ws", WebSocketHandler(hub))

	fmt.Println("🌐 Web UI → http://localhost:8080")
	fmt.Println("🔌 WS → ws://localhost:8080/ws")
	fmt.Println("🔌 Ring → ws://localhost:8080/ring")

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", Port), nil))
}
