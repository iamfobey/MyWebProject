package handlers

import (
	"MyWebProject/application"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var jwtKey = []byte("your_secret_key")

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func generateJWT(login string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthHandler(application *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
			return
		}

		if r.FormValue("authType") == "signup" {
			var creds Credentials

			creds.Login = r.FormValue("login")
			creds.Password = r.FormValue("password")

			if creds.Login == "" || creds.Password == "" {
				http.Error(w, "Недопустимые данные для регистрации", http.StatusBadRequest)
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
				return
			}

			_, err = application.Database.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", creds.Login, hashedPassword)
			if err != nil {
				var err *pq.Error

				if errors.As(err, &err) && err.Code == "23505" {
					http.Error(w, "Логин уже существует", http.StatusConflict)
				}
				return
			}

			token, err := generateJWT(creds.Login)
			if err != nil {
				http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(Response{Message: "Успешная регистрация", Token: token})

			fmt.Println("Новый пользователь:", creds.Login)
		} else if r.FormValue("authType") == "login" {
			var creds Credentials
			creds.Login = r.FormValue("login")
			creds.Password = r.FormValue("password")

			var storedPassword string
			err := application.Database.QueryRow("SELECT password FROM users WHERE login=$1", creds.Login).Scan(&storedPassword)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
				} else {
					http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
				}
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
			if err != nil {
				http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
				return
			}

			token, err := generateJWT(creds.Login)
			if err != nil {
				http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(Response{Message: "Успешный вход", Token: token})
		}
	}
}
