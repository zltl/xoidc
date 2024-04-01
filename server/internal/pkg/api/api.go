package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/internal/pkg/storage"
)

type Handler struct {
	Store *storage.Storage
}

// serve /api/oidc/...
func (h *Handler) Serve(r chi.Router) {
	r.Get("/clients", h.handleGetClientList)
	r.Post("/clients", h.handlePostClient)
	r.Get("/clients/{client_id}", h.handleGetClient)
	// r.Get("/", h.index)
}

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

func (h *Handler) decodeJSON(ctx context.Context, r *http.Request, v any) error {
	_ = ctx
	// read all data
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	logrus.Infof("url=%s body=%s", r.URL, buf)
	err = json.Unmarshal(buf, &v)
	return err
}
