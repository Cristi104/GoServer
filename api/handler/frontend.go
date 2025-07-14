package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func isDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func FrontendHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.URL.String(), "/Messanger") == 0 {
		_, err := authentifacateUser(r)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	buildDir := "./frontend/dist"
	path := filepath.Join(buildDir, r.URL.Path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) || isDirectory(path) {
		http.ServeFile(w, r, filepath.Join(buildDir, "index.html"))
		return
	}
	fs := http.FileServer(http.Dir("frontend/dist/"))
	fs.ServeHTTP(w, r)
}
