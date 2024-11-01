package main

import (
	"MyWebProject/application"
	"MyWebProject/handlers"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	app := application.NewApplication()

	mux := http.NewServeMux()

	handler := cors.Default().Handler(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler = c.Handler(handler)

	mux.HandleFunc("/api/public/auth/", handlers.AuthHandler(app))

	log.Fatal(http.ListenAndServe(":8080", handler))
}
