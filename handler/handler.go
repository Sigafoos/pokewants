package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sigafoos/pokewants/wants"
)

type Handler struct {
	want *wants.Wants
}

type Request struct {
	User    string `json:"user"`
	Pokemon string `json:"pokemon"`
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
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
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

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading POST body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req Request
	err = json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("error unmarshalling POST request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.User == "" || req.Pokemon == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.want.Add(req.User, req.Pokemon)
	if err != nil {
		if err == wants.ErrorPokemonNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("error adding want: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading DELETE body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req Request
	err = json.Unmarshal(b, &req)
	if err != nil {
		log.Printf("error unmarshalling DELETE request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.User == "" || req.Pokemon == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.want.Delete(req.User, req.Pokemon)
	if err != nil {
		if err == wants.ErrorPokemonNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("error deleting want: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
