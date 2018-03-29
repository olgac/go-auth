package main

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
)

func Health(w http.ResponseWriter, r *http.Request) {
	responseOK(w, Status{Status: "UP"})
}

func Token(w http.ResponseWriter, r *http.Request) {
	apikey := r.Header.Get(APIKEY_HEADER)

	if len(apikey) == 0 {
		response(w, http.StatusBadRequest, Response{Code: http.StatusBadRequest, Text: "APIKEY header is required!"})
		return
	}

	if val, ok := APIKEYS["apikey:"+apikey]; ok {

		secret := []byte(AUTH_SECRET)

		issuer := jwt.New(jwt.SigningMethodHS256)
		issuer.Claims = jwt.MapClaims{
			"sub":   val.Name,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
			"roles": strings.Split(val.Roles, ","),
		}

		token, err := issuer.SignedString(secret)
		if err != nil {
			responseError(w, err)
			return
		}

		responseOK(w, JWT{Token: token})
	} else {
		response(w, http.StatusUnauthorized, Response{Code: http.StatusUnauthorized, Text: "APIKEY is unauthorized!"})
	}
}

func Cache(w http.ResponseWriter, r *http.Request) {
	c, err := redis.Dial("tcp", REDIS_HOST+":"+REDIS_PORT)
	if err != nil {
		responseError(w, err)
		return
	}
	defer c.Close()

	keys, err := redis.Strings(c.Do("KEYS", "apikey:*"))
	if err != nil {
		responseError(w, err)
		return
	}

	APIKEYS = make(map[string]ApiKey)

	var a ApiKey

	for _, key := range keys {
		v, err := redis.Values(c.Do("HGETALL", key))
		if err != nil {
			responseError(w, err)
			return
		}

		err2 := redis.ScanStruct(v, &a)
		if err2 != nil {
			responseError(w, err)
			return
		}
		APIKEYS[key] = a
	}

	responseOK(w, Response{Code: http.StatusOK, Text: "OK"})
}
