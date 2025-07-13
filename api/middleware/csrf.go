package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/csrf"
)

var CsrfMiddleware func(http.Handler) http.Handler

func init() {
	type csrfJson struct {
		Key string
		Secure bool
		TrustedOrigins []string
	}

	token := csrfJson{"default_token_for_csrf_protectio", false, []string{"localhost:8080"}}

	readData, err := os.ReadFile("config/csrf.json")
	if err != nil {
		data, err1 := json.Marshal(token)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.MkdirAll(filepath.Join(".", "config"), os.ModePerm)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.WriteFile("config/csrf.json", data, 0600)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Println("Missing configuration file created, please change the CSRF key information inside config/csrf.json")
		
	} else {
		err = json.Unmarshal(readData, &token)
		if err != nil {
			log.Fatal(err)
		}

		if token.Key == "default_token_for_csrf_protectio" {
			log.Println("Warning: CSRF key is set to the default value, please change it.")
		}

	}

	CsrfMiddleware = csrf.Protect(
		[]byte(token.Key), 
		csrf.Secure(token.Secure), 
		csrf.HttpOnly(true), 
		csrf.Path("/"), 
		csrf.CookieName("csrf"),
		csrf.TrustedOrigins([]string{"localhost:8080", "192.168.0.137:8080"}),
	)
}
