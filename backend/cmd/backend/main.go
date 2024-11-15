package main

import (
	httpadapter "backend/internal/backend/adapters/http-adapter"
	"backend/internal/backend/adapters/http-adapter/handlers"
	"backend/internal/backend/adapters/http-adapter/router"
	postgresadapter "backend/internal/backend/adapters/postgres-adapter"
	"backend/internal/backend/application"
	userstorepostgres "backend/internal/backend/store/user-store-postgres"
	"log"
	"net/http"
	"os"
)

func main() {
	httpAdapter := httpadapter.New(
		httpadapter.Config{
			Port: "8080",

			AllowedOrigins:   []string{"http://192.168.0.106:3000"},
			AllowCredentials: true,
			AllowedMethods:   []string{"POST", "GET"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			Debug:            true,
		},
	)

	postgresAdapter := postgresadapter.New(
		postgresadapter.Config{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DbName:   os.Getenv("POSTGRES_DB"),
		},
	)

	httpAdapter.Router.AddRoute(router.Route{
		Pattern: "/api/public/auth/",
		Handler: &handlers.AuthHandler{
			UserStore: userstorepostgres.New(userstorepostgres.Config{
				Database:      postgresAdapter.Database,
				AdminLogin:    os.Getenv("USER_STORE_ADMIN_LOGIN"),
				AdminPassword: os.Getenv("USER_STORE_ADMIN_PASSWORD"),
			}),
			JWTKey: os.Getenv("JWT_KEY"),
		},
		HTTPMethod: http.MethodPost,
	})

	app := application.New(application.Config{
		HttpAdapter:     httpAdapter,
		PostgresAdapter: postgresAdapter,
	})

	if err := app.Run(); err != nil {
		log.Printf("Run error: %s", err.Error())
	}

	defer func(app *application.App) {
		if err := app.Close(); err != nil {
			log.Printf("Close error: %s", err.Error())
		}
	}(app)
}
