package main

import (
	"encoding/json"
	"net/http"

	"github.com/garyburd/redigo/redis"
)

func Health(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	var response = Status{Status: "UP"}
	errorHandler(json.NewEncoder(w).Encode(response), "Health")
}

func Token(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	var response = JWT{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG54eHggRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.9cZtc89mhRfFSXMcGHb4lKsD8gCMACAduMilx0oJD2E"}
	errorHandler(json.NewEncoder(w).Encode(response), "Token")
}

func Cache(w http.ResponseWriter, r *http.Request) {
	c, err := redis.Dial("tcp", REDIS_HOST+":"+REDIS_PORT)
	defer c.Close()
	errorHandler(err, "Redis Dial")

	keys, err := redis.Strings(c.Do("KEYS", "apikey:*"))
	errorHandler(err, "Redis KEYS")

	APIKEYS := make(map[string]ApiKey)

	var a ApiKey

	for _, key := range keys {
		//fmt.Println(key)
		v, err2 := redis.Values(c.Do("HGETALL", key))
		errorHandler(err2, "Redis HGETALL")

		err := redis.ScanStruct(v, &a)
		errorHandler(err, "Redis ScanStruct")
		APIKEYS[key] = a
		//fmt.Printf("%+v\n", APIKEYS[key])
	}

	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	var response = Response{Code: http.StatusOK, Text: "OK"}
	errorHandler(json.NewEncoder(w).Encode(response), "APIKey")
}

func setHeaders(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return w
}
