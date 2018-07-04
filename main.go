package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	LOG_LEVEL     = os.Getenv("LOG_LEVEL")
	AUTH_SECRET   = os.Getenv("AUTH_SECRET")
	APP_PORT      = os.Getenv("APP_PORT")
	REDIS_HOST    = os.Getenv("REDIS_HOST")
	REDIS_PORT    = os.Getenv("REDIS_PORT")
	APIKEYS       = make(map[string]ApiKey)
	APIKEY_HEADER = "apikey"
)

func init() {
	if LOG_LEVEL == "" {
		LOG_LEVEL = "warn"
	}
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
	if err := http.ListenAndServe(":"+APP_PORT, router); err != nil {
		log.Fatalln(err)
	}
}

func response(w http.ResponseWriter, code int, text string) {
	responseObject(w, code, Response{Code: code, Text: text})
}

func responseObject(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
