package http_adapter

import (
	"backend/internal/backend/adapters/http-adapter/router"
	"github.com/rs/cors"
	"net/http"
)

type Adapter struct {
	ServeMux *http.ServeMux
	Handler  http.Handler
	Router   *router.Router
	Config   Config
}

type Config struct {
	Port string

	AllowedOrigins   []string
	AllowCredentials bool
	AllowedMethods   []string
	AllowedHeaders   []string
	Debug            bool
}

func New(config Config) *Adapter {
	mux := http.NewServeMux()

	handler := cors.Default().Handler(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowCredentials: config.AllowCredentials,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		Debug:            config.Debug,
	})

	handler = c.Handler(handler)

	routerNew := router.New(router.Config{})

	return &Adapter{
		ServeMux: mux,
		Handler:  handler,
		Router:   routerNew,
		Config:   config,
	}
}

func (adapter *Adapter) Run() error {
	for _, route := range adapter.Router.Routes {
		adapter.ServeMux.Handle(route.Pattern, route.Handler)
	}

	err := http.ListenAndServe(":"+adapter.Config.Port, adapter.Handler)

	if err != nil {
		return err
	}
	return nil
}
