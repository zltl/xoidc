package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/internal/pkg/storage"
)

type Handler struct {
	Store *storage.Storage
}

// serve /api/...
func (h *Handler) Serve(r chi.Router) {
	r.Get("/client", h.handleGetClientList)
	// r.Get("/", h.index)
}

type Response struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

const (
	Success          = "success"
	ErrFailed        = "failed"
	ErrInvalidParams = "invalid_params"
)

func (h *Handler) R(w http.ResponseWriter, r *http.Request, code int, res any) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"failed","msg":"internal server error"}`))
	}
	w.WriteHeader(code)
	w.Write(body)
	logrus.Infof("%s %s -> %s", r.Method, r.URL.String(), body)
}
