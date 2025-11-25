package main

import (
	"encoding/json"
	"net/http"
)

func WebSocketHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub.Listen(w, r)
	}
}

type RingBody struct {
	Password string `json:"password"`
}

func RingHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)

		body := &RingBody{}
		err := decoder.Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if body.Password != Password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hub.Broadcast("ring")
	}
}
