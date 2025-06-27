package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JsonPass обёртка над password для удобства вывода в формате json.
type JsonPass struct {
	Password string `json:"password"`
}

// JsonPass обёртка над token для удобства вывода в формате json.
type JsonToken struct {
	Token string `json:"token"`
}

var envPass = os.Getenv("TODO_PASSWORD") // Получаем переменную окружения TODO_PASSWORD.

func passCheckHandler(w http.ResponseWriter, r *http.Request) {
	var pass JsonPass
	var token JsonToken
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	if len(envPass) > 0 {
		if envPass != pass.Password {
			w.WriteHeader(http.StatusUnauthorized)
			writeJsonErr(w, fmt.Errorf("неверный пароль"))
			return
		}
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		token.Token, err = jwtToken.SignedString([]byte(envPass))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeJsonErr(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJson(w, token)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(envPass) > 0 {
			var jwtFromCookie string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtFromCookie = cookie.Value
			}
			jwtNew := jwt.New(jwt.SigningMethodHS256)
			jwtFromPass, err := jwtNew.SignedString([]byte(envPass))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				writeJsonErr(w, err)
				return
			}
			if jwtFromCookie != jwtFromPass {
				w.WriteHeader(http.StatusUnauthorized)
				writeJsonErr(w, fmt.Errorf("требуется аутентификация"))
				return
			}
		}
		next(w, r)
	})
}
