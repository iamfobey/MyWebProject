package router

import "net/http"

type Config struct {
}

type Router struct {
	Config Config

	Routes Routes
}

func New(config Config) *Router {
	return &Router{
		Config: config,
	}
}

func (router *Router) AddRoute(route Route) {
	router.Routes = append(router.Routes, route)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range router.Routes {
		if req.URL.Path == route.Pattern && req.Method == route.HTTPMethod {
			route.Handler.ServeHTTP(w, req)
			return
		}
	}

	http.NotFound(w, req)
}
