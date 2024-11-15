package test

import (
	httpadapter "backend/internal/backend/adapters/http-adapter"
	"backend/internal/backend/adapters/http-adapter/handlers"
	"backend/internal/backend/adapters/http-adapter/router"
	postgresadapter "backend/internal/backend/adapters/postgres-adapter"
	"backend/internal/backend/application"
	userstorepostgres "backend/internal/backend/store/user-store-postgres"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func FuzzAuthHandler(f *testing.F) {
	f.Add("signup", "testuser", "Test@123456")
	f.Add("login", "testuser", "Test@123456")
	f.Add("signup", "", "")

	f.Fuzz(func(t *testing.T, authType, login, password string) {
		httpAdapter := httpadapter.New(httpadapter.Config{
			Port:             "8080",
			AllowedOrigins:   []string{"http://192.168.0.106:3000"},
			AllowCredentials: true,
			AllowedMethods:   []string{"POST", "GET"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			Debug:            true,
		})

		postgresAdapter := postgresadapter.New(postgresadapter.Config{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DbName:   os.Getenv("POSTGRES_DB"),
		})

		httpAdapter.Router.AddRoute(router.Route{
			Pattern: "/auth",
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

		application.New(application.Config{
			HttpAdapter:     httpAdapter,
			PostgresAdapter: postgresAdapter,
		})

		formData := "authType=" + authType + "&login=" + login + "&password=" + password
		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBufferString(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		httpAdapter.Router.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if authType == "signup" && resp.StatusCode == http.StatusConflict {
			t.Log("signup StatusConflict")
		}

		if authType == "signup" && password != "" && login != "" && resp.StatusCode != http.StatusOK {
			t.Errorf("signup StatusCode != StatusOK")
		}

		if authType == "login" && password != "" && resp.StatusCode == http.StatusUnauthorized {
			t.Log("login StatusUnauthorized")
		}
	})
}
