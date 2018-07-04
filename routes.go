package main

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
