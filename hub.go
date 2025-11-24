package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]bool)}
}

func (h *Hub) Add(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	log.Println("Client connected:", conn.RemoteAddr())
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()

	log.Println("Client removed:", conn.RemoteAddr())
}

func (h *Hub) Broadcast(message string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Client send error:", err)

			conn.Close()
			delete(h.clients, conn)
		}
	}
}

func (h *Hub) Listen(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade (listen) failed:", err)
		return
	}

	h.Add(conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.Remove(conn)

			conn.Close()
			return
		}
	}
}

func (h *Hub) Control(w http.ResponseWriter, r *http.Request) {
	pw := r.URL.Query().Get("password")
	if pw != Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade (control) failed:", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		switch string(msg) {
		case "ring":
			h.Broadcast("ring")
		default:
			log.Println("Unknown control command:", string(msg))
		}
	}
}
