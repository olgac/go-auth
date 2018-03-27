package main

import "net/http"

type Response struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type ApiKey struct {
	Name   string
	Secret string
	Claims string
}

type Status struct {
	Status string `json:"status"`
}

type JWT struct {
	Token string `json:"token"`
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/health",
		Health,
	},
	Route{
		"GET",
		"/token",
		Token,
	},
	Route{
		"GET",
		"/cache",
		Cache,
	},
}
