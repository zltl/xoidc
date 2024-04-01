package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zltl/xoidc/server/internal/pkg/m"

	"github.com/sirupsen/logrus"
)

// get one client by id
// GET /api/oidc/clients/{clientid}
func (h *Handler) handleGetClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clientID := chi.URLParam(r, "client_id")
	cli, err := h.Store.GetClient(ctx, clientID)
	if err != nil {
		logrus.Error(err)
		h.R(w, r, http.StatusInternalServerError, m.Response{
			Status: m.ErrFailed,
			Msg:    err.Error(),
		})
		return
	}
	h.R(w, r, http.StatusOK, m.ClientResponse{
		Response: m.Response{
			Status: m.Success,
		},
		Client: m.ClientDB2View(cli), // TODO: fix it
	})
}

// list all clients
// GET /api/oidc/clients?limit=10&offset=0
func (h *Handler) handleGetClientList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	total, err := h.Store.TotalClient(ctx)
	if err != nil {
		h.R(w, r, http.StatusInternalServerError, m.Response{
			Status: m.ErrFailed,
			Msg:    err.Error(),
		})
		return
	}

	clients, err := h.Store.GetAllClient(ctx, offset, limit)
	if err != nil {
		h.R(w, r, http.StatusInternalServerError, m.Response{
			Status: m.ErrFailed,
			Msg:    err.Error(),
		})
		return
	}

	cs := []m.Client{}
	for _, c := range clients {
		cs = append(cs, m.ClientDB2View(&c))
	}

	h.R(w, r, http.StatusOK, m.ClientListResponse{
		Response: m.Response{
			Status: m.Success,
			Msg:    "success",
		},
		Total:   total,
		Clients: cs,
	})
}

// POST /api/oidc/client
// create new client
func (h *Handler) handlePostClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var client m.Client
	// parse json
	if err := h.decodeJSON(ctx, r, &client); err != nil {
		logrus.Error(err)
		h.R(w, r, http.StatusBadRequest, m.Response{
			Status: m.ErrInvalidRequest,
			Msg:    err.Error(),
		})
	}
	// TODO: save to db
	// TODO: return client to client

}
