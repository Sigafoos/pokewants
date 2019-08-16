package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sigafoos/pokewants/wants"
)

type Handler struct {
	want *wants.Wants
}

func New(want *wants.Wants) *Handler {
	return &Handler{
		want: want,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := h.want.Get(user)

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}
