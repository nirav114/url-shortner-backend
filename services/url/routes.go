package url

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nirav114/url-shortner-backend.git/types"
)

type Handler struct {
	store types.UrlStore
}

func NewHandler(store types.UrlStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/saveUrl", h.handleSaveUrl).Methods("POST")
	router.HandleFunc("/modifyUrl", h.handleModifyUrl).Methods("POST")
	router.HandleFunc("/removeUrl", h.handleRemoveUrl).Methods("POST")
	router.HandleFunc("/getAllUrls", h.handleGetAllUrls).Methods("POST")
}

func (h *Handler) handleSaveUrl(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleModifyUrl(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRemoveUrl(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetAllUrls(w http.ResponseWriter, r *http.Request) {

}
