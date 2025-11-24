package main

import (
	"net/http"
)

func WebSocketHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("password") != "" {
			hub.Control(w, r)
		} else {
			hub.Listen(w, r)
		}
	}
}
