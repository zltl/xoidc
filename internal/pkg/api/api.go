package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/zltl/xoidc/internal/pkg/storage"
)

type Handler struct {
	Store *storage.Storage
}

// serve /api/...
func (h *Handler) Serve(r chi.Router) {
	r.Get("/client", h.handleClient)
	// r.Get("/", h.index)
}
