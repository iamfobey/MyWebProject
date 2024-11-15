package router

import "net/http"

type HTTPMethod string

type Route struct {
	Pattern    string
	Handler    http.Handler
	HTTPMethod string
}

type Routes []Route
