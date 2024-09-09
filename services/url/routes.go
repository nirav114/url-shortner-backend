package url

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nirav114/url-shortner-backend.git/types"
	"github.com/nirav114/url-shortner-backend.git/utils"
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
	var payload types.SaveUrlPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		err = err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err := h.store.GetUrlByShortUrl(payload.ShortUrl)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("url with this short url /%s already exists", payload.ShortUrl))
		return
	}

	err = h.store.CreateUrl(types.Url{
		ShortUrl: payload.ShortUrl,
		FullUrl:  payload.FullUrl,
		UserID:   payload.UserID,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleModifyUrl(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRemoveUrl(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetAllUrls(w http.ResponseWriter, r *http.Request) {

}
