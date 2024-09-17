package url

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mssola/user_agent"
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
	router.Handle("/getStats", auth.JWTMiddleware(userStore, http.HandlerFunc(h.handleGetStats))).Methods("POST")
	router.HandleFunc("/r/{shortUrl}", h.handleURLRedirect).Methods("GET")
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

	userClaims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}
	userID := userClaims.UserID

	oldUrl, err := h.store.GetUrlByShortUrl(payload.ShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("url with this short url /%s doesn't exist", payload.ShortUrl))
		return
	}

	if oldUrl.UserID != userID {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid url modification"))
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

	userClaims, ok := auth.GetUserClaimsFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}
	userID := userClaims.UserID

	url, err := h.store.GetUrlByShortUrl(payload.ShortUrl)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("url with this short url /%s doesn't exist", payload.ShortUrl))
		return
	}

	if url.UserID != userID {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid url removal request"))
		return
	}

	err = h.store.RemoveUrl(payload.ShortUrl)
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

func (h *Handler) handleURLRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"]

	url, err := h.store.GetUrlByShortUrl(shortUrl)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	ip := utils.GetIPAddress(r)
	ua := user_agent.New(r.UserAgent())
	browser, _ := utils.GetBrowserInfo(ua)
	device := utils.GetDeviceType(ua)
	platform := utils.GetPlatform(ua)
	language := utils.GetLanguage(r)
	country := utils.GetCountryFromIP(ip)

	err = h.store.InsertClickData(url.ID, ip, country, device, platform, browser, language)
	if err != nil {
		http.Error(w, "Failed to track click", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.FullUrl, http.StatusMovedPermanently)
}

func (h *Handler) handleGetStats(w http.ResponseWriter, r *http.Request) {
	var payload types.GetStatsPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		err = err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	url, err := h.store.GetUrlByShortUrl(payload.ShortUrl)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	clicks, err := h.store.GetClicksByID(url.ID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	countryCount := make(map[string]int)
	deviceCount := make(map[string]int)
	platformCount := make(map[string]int)
	browserCount := make(map[string]int)
	languageCount := make(map[string]int)
	ipClicks := make(map[string]bool)

	for _, click := range clicks {
		countryCount[click.Country]++
		deviceCount[click.Device]++
		platformCount[click.Platform]++
		browserCount[click.Browser]++
		languageCount[utils.GetPrimaryLanguage(click.Language)]++
		ipClicks[click.IPAddress] = true
	}

	clickType := map[string]int{
		"unique": len(ipClicks),
		"total":  len(clicks),
	}

	stats := map[string]interface{}{
		"country_count":  countryCount,
		"device_count":   deviceCount,
		"platform_count": platformCount,
		"browser_count":  browserCount,
		"language_count": languageCount,
		"click_type":     clickType,
	}

	utils.WriteJSON(w, http.StatusOK, stats)
}
