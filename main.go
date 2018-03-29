package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var (
	AUTH_SECRET   = os.Getenv("AUTH_SECRET")
	APP_PORT      = os.Getenv("APP_PORT")
	REDIS_HOST    = os.Getenv("REDIS_HOST")
	REDIS_PORT    = os.Getenv("REDIS_PORT")
	APIKEYS       = make(map[string]ApiKey)
	APIKEY_HEADER = "apikey"
)

func init() {
	if AUTH_SECRET == "" {
		AUTH_SECRET = "AhmetVehbiOlgac"
	}
	if APP_PORT == "" {
		APP_PORT = "8080"
	}
	if REDIS_HOST == "" {
		REDIS_HOST = "localhost"
	}
	if REDIS_PORT == "" {
		REDIS_PORT = "6379"
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Handler(route.HandlerFunc)
	}
	s := &http.Server{
		Handler:      router,
		Addr:         ":" + APP_PORT,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Panic(s.ListenAndServe())
}

func responseError(w http.ResponseWriter, err error) {
	defer log.Printf("%s - %v", "ERROR", err)
	response(w, http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Text: err.Error()})
}

func responseOK(w http.ResponseWriter, v interface{}) {
	response(w, http.StatusOK, v)
}

func response(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
