package handler

import (
	"GoServer/repository"
	"fmt"
	"net/http"
	"time"
	"encoding/json"

	jwt "github.com/golang-jwt/jwt/v5"
)

const SERCURE = false;

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

func authentifacateUser(r *http.Request) (repository.User, error){
	auth, err := r.Cookie("auth")
	if err != nil {
		return repository.User{}, err
	}

	token, err := validateJWT(auth.Value)
	if err != nil {
		return repository.User{}, err
	}

	var user repository.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp := claims["exp"].(float64)
		expTime := time.Unix(int64(exp), 0)
		if expTime.Before(time.Now()) {
			return repository.User{}, err
		}

		userData := claims["userData"].(string)
		err = json.Unmarshal([]byte(userData), &user)
		if err != nil {
			return repository.User{}, err
		}
	} else {
		return repository.User{}, err
	}

	return user, nil
}
