package main

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
)

func Health(w http.ResponseWriter, r *http.Request) {
	response(w, http.StatusOK, "UP")
}

func Token(w http.ResponseWriter, r *http.Request) {
	if apikey := r.Header.Get(APIKEY_HEADER); len(apikey) == 0 {
		response(w, http.StatusBadRequest, "APIKEY header is required!")
	} else if val, ok := APIKEYS["apikey:"+apikey]; ok {
		secret := []byte(AUTH_SECRET)
		issuer := jwt.New(jwt.SigningMethodHS256)
		issuer.Claims = jwt.MapClaims{
			"sub":   val.Name,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
			"roles": strings.Split(val.Roles, ","),
		}

		if token, err := issuer.SignedString(secret); err != nil {
			response(w, http.StatusInternalServerError, err.Error())
		} else {
			responseObject(w, http.StatusOK, JWT{Token: token})
		}
	} else {
		response(w, http.StatusUnauthorized, "APIKEY is unauthorized!")
	}
}

func Cache(w http.ResponseWriter, r *http.Request) {
	if c, err := redis.Dial("tcp", REDIS_HOST+":"+REDIS_PORT); err != nil {
		response(w, http.StatusInternalServerError, err.Error())
	} else {
		defer c.Close()
		if keys, err := redis.Strings(c.Do("KEYS", "apikey:*")); err != nil {
			response(w, http.StatusInternalServerError, err.Error())
		} else {
			APIKEYS = make(map[string]ApiKey)
			var a ApiKey
			for _, key := range keys {
				if v, err := redis.Values(c.Do("HGETALL", key)); err != nil {
					response(w, http.StatusInternalServerError, "key : "+key+" "+err.Error())
					return
				} else if err := redis.ScanStruct(v, &a); err != nil {
					response(w, http.StatusInternalServerError, err.Error())
					return
				} else {
					APIKEYS[key] = a
				}
			}
			response(w, http.StatusOK, "OK")
		}
	}
}
