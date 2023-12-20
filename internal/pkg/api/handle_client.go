package api

import "net/http"

// list all clients
func (h *Handler) handleClient(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
