package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/mssola/user_agent"
	"github.com/nirav114/url-shortner-backend.git/config"
	"github.com/oschwald/geoip2-golang"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func GetIPAddress(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func GetDeviceType(ua *user_agent.UserAgent) string {
	if ua.Mobile() {
		return "mobile"
	}
	return "desktop"
}

func GetBrowserInfo(ua *user_agent.UserAgent) (string, string) {
	return ua.Browser()
}

func GetPlatform(ua *user_agent.UserAgent) string {
	return ua.Platform()
}

func GetCountryFromIP(ip string) string {
	db, err := geoip2.Open(config.EnvConfig.GEO_LITE_DB)
	if err != nil {
		log.Println("Failed to open GeoIP2 database:", err)
		return "unknown"
	}
	defer db.Close()

	ipAddr := net.ParseIP(ip)

	if ipAddr == nil {
		log.Println("Invalid IP address:", ip)
		return "unknown"
	}

	record, err := db.Country(ipAddr)
	if err != nil {
		log.Println("Failed to lookup country for IP:", err)
		return "unknown"
	}

	if countryName, ok := record.Country.Names["en"]; ok {
		return countryName
	}
	return "unknown"
}

func GetLanguage(r *http.Request) string {
	return r.Header.Get("Accept-Language")
}
