package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

func isDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func FrontendHandler(w http.ResponseWriter, r *http.Request) {
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
