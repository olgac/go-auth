package main

import (
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
		router.Methods(route.Method).Path(route.Pattern).Handler(withLogger(route.HandlerFunc))
	}
	s := &http.Server{
		Handler:      router,
		Addr:         ":" + APP_PORT,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	errorHandler(s.ListenAndServe(), "Server")
}

func withLogger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf("%s - %s", "INFO", r.Method+" "+r.RequestURI)
		inner.ServeHTTP(w, r)
	})
}

func errorHandler(err error, str string) {
	if err != nil {
		log.Fatalf("%s - %v - %q", "ERROR", err, str)
	}
}
