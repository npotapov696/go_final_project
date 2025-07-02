package tests

import (
	"os"

	"github.com/golang-jwt/jwt"
)

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = getToken()

// getToken получает заначение токена, исходя из пароля, лежащего в переменной среды окружения TODO_PASSWORD.
func getToken() string {
	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) > 0 {
		jwtNew := jwt.New(jwt.SigningMethodHS256)
		jwtFromPass, err := jwtNew.SignedString([]byte(pass))
		if err != nil {
			return ``
		}
		return jwtFromPass
	}
	return ``
}
