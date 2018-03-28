package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	var response Response
	apikey := r.Header.Get(APIKEY_HEADER)
	fmt.Println("1 : " + apikey)

	if len(apikey) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response = Response{Code: http.StatusBadRequest, Text: "APIKEY header is required!"}
		errorHandler(json.NewEncoder(w).Encode(response), "APIKEY")
		return
	}
	for k := range APIKEYS {
		fmt.Printf("key[%s] value[%s]\n", k, APIKEYS[k])
	}

	if val, exists := APIKEYS["apikey:"+apikey]; exists {

		secret := []byte(AUTH_SECRET)

		fmt.Println("2 : " + val.Name + " " + val.Roles)

		issuer := jwt.New(jwt.SigningMethodHS256)
		issuer.Claims = jwt.MapClaims{
			"sub":   val.Name,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
			"roles": val.Roles,
		}

		token, err := issuer.SignedString(secret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("%s - %v - %q", "ERROR", err, "Token SignedString")
			response = Response{Code: http.StatusInternalServerError, Text: "APIKEY signing error!"}
			errorHandler(json.NewEncoder(w).Encode(response), "Token SignedString")
			return
		}

		w.WriteHeader(http.StatusOK)
		errorHandler(json.NewEncoder(w).Encode(JWT{Token: token}), "Token SignedString")
	} else {

		fmt.Println("2 : " + val.Name)
		w.WriteHeader(http.StatusUnauthorized)
		response = Response{Code: http.StatusUnauthorized, Text: "APIKEY is unauthorized!"}
		errorHandler(json.NewEncoder(w).Encode(response), "APIKEY")
	}
}

func Cache(w http.ResponseWriter, r *http.Request) {
	c, err := redis.Dial("tcp", REDIS_HOST+":"+REDIS_PORT)
	defer c.Close()
	errorHandler(err, "Redis Dial")

	keys, err := redis.Strings(c.Do("KEYS", "apikey:*"))
	errorHandler(err, "Redis KEYS")

	APIKEYS = make(map[string]ApiKey)

	var a ApiKey

	for _, key := range keys {
		//fmt.Println(key)
		v, err2 := redis.Values(c.Do("HGETALL", key))
		errorHandler(err2, "Redis HGETALL")

		err := redis.ScanStruct(v, &a)
		errorHandler(err, "Redis ScanStruct")
		APIKEYS[key] = a
		fmt.Printf("%+v\n", APIKEYS[key])
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
