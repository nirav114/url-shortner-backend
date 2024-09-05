package user

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("in the login")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("in the register")
	w.WriteHeader(http.StatusOK)
}
