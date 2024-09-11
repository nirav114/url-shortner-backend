package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nirav114/url-shortner-backend.git/services/url"
	"github.com/nirav114/url-shortner-backend.git/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr, db}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	urlStore := url.NewStore(s.db)
	urlHandler := url.NewHandler(urlStore)
	urlHandler.RegisterRoutes(userStore, subrouter)

	log.Println("Serving on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
