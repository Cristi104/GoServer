package handler

import (
	"net/http"
	"time"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

func createAuthJWT(userData string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userData": userData,
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	})

	signedToken, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTKey), nil
	}), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func errorResponseJson(w http.ResponseWriter, err string, code int) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"success\":%s, \"error\":\"%s\"}", "false", err)))
}
