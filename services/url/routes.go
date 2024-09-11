package url

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nirav114/url-shortner-backend.git/services/auth"
	"github.com/nirav114/url-shortner-backend.git/types"
	"github.com/nirav114/url-shortner-backend.git/utils"
)

type Handler struct {
	store types.UrlStore
}

func NewHandler(store types.UrlStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(userStore types.UserStore, router *mux.Router) {
	router.Handle("/saveUrl", auth.JWTMiddleware(userStore, http.HandlerFunc(h.handleSaveUrl))).Methods("POST")
	router.Handle("/modifyUrl", auth.JWTMiddleware(userStore, http.HandlerFunc(h.handleModifyUrl))).Methods("POST")
	router.Handle("/removeUrl", auth.JWTMiddleware(userStore, http.HandlerFunc(h.handleRemoveUrl))).Methods("POST")
	router.Handle("/getAllUrls", auth.JWTMiddleware(userStore, http.HandlerFunc(h.handleGetAllUrls))).Methods("POST")
}

func (h *Handler) handleSaveUrl(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := auth.GetUserClaimsFromContext(r.Context())

	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}
	userID := userClaims.UserID

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
		UserID:   userID,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleModifyUrl(w http.ResponseWriter, r *http.Request) {
	var payload types.ModifyUrlPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		err = err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	oldUrl, err := h.store.GetUrlByShortUrl(payload.ShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("url with this short url /%s doesn't exist", payload.ShortUrl))
		return
	}

	newUrl := oldUrl
	newUrl.FullUrl = payload.FullUrl
	err = h.store.ModifyUrl(*oldUrl, *newUrl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleRemoveUrl(w http.ResponseWriter, r *http.Request) {
	var payload types.RemoveUrlPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		err = err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.store.RemoveUrl(payload.ShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleGetAllUrls(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}
	userID := userClaims.UserID

	urls, err := h.store.GetUrlsByUserID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	log.Println(urls)
	utils.WriteJSON(w, http.StatusOK, map[string][]*types.UrlResponse{"urls": urls})
}
