package handlers

import (
	"backend/internal/backend/domain/user"
	jwtutils "backend/pkg/backend/jwt-utils"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// Response TODO: may be to Route struct?
type Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func LogHttpError(w http.ResponseWriter, httpMessage string, httpStatus int, logMessage string) {
	log.Printf("HTTP error '%d': %s - %s \n", httpStatus, httpMessage, logMessage)
	http.Error(w, httpMessage, httpStatus)
}

type AuthHandler struct {
	UserStore user.Store
	JWTKey    string
}

func (handler *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		LogHttpError(w, "Метод не поддерживается", http.StatusMethodNotAllowed, fmt.Sprintf("%v", r.Method))
		return
	}

	if err := r.ParseForm(); err != nil {
		LogHttpError(w, "Ошибка при разборе формы", http.StatusBadRequest, fmt.Sprintf("%v", r.Method))
		return
	}

	if r.FormValue("authType") == "signup" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost+4)
		if err != nil {
			LogHttpError(w, "Ошибка при хешировании пароля", http.StatusInternalServerError, fmt.Sprintf("%v", r.Method))
			return
		}

		var usr = user.NewUser(user.Credentials{
			Login:    r.FormValue("login"),
			Password: string(hashedPassword),
		})

		err = handler.UserStore.CreateUser(usr)

		if err != nil {
			var pqErr *pq.Error

			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				LogHttpError(w, "Логин уже существует", http.StatusConflict, fmt.Sprintf("%v", r.Method))
			}
			return
		}

		token, err := jwtutils.GenerateFromLogin(usr.Credentials.Login, handler.JWTKey)
		if err != nil {
			LogHttpError(w, "Ошибка генерации токена", http.StatusInternalServerError, fmt.Sprintf("%v", r.Method))
			return
		}

		json.NewEncoder(w).Encode(Response{Message: "Успешная регистрация", Token: token})
	} else if r.FormValue("authType") == "login" {
		usr, err := handler.UserStore.GetUserByLogin(r.FormValue("login"))

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				LogHttpError(w, "Неверный логин или пароль", http.StatusUnauthorized, fmt.Sprintf("%v", r.Method))
			} else {
				LogHttpError(w, "Ошибка базы данных", http.StatusInternalServerError, fmt.Sprintf("%v", r.Method))
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(usr.Credentials.Password), []byte(r.FormValue("password")))

		if err != nil {
			LogHttpError(w, "Неверный логин или пароль", http.StatusUnauthorized, fmt.Sprintf("%v", r.Method))
			return
		}

		token, err := jwtutils.GenerateFromLogin(usr.Credentials.Login, handler.JWTKey)
		if err != nil {
			LogHttpError(w, "Ошибка генерации токена", http.StatusInternalServerError, fmt.Sprintf("%v", r.Method))
			return
		}

		json.NewEncoder(w).Encode(Response{Message: "Успешный вход", Token: token})
	}
}
