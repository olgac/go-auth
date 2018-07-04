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
	Name  string
	Roles string
}

type JWT struct {
	Token string `json:"token"`
}

type Level int

const (
	DEBUG Level = 1 + iota
	INFO
	WARN
	ERROR
	FATAL
)
