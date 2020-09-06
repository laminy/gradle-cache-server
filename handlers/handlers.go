package handlers

import (
	"github.com/laminy/gradle-cache-server/config"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func CachePut(w http.ResponseWriter, r *http.Request) {
	fn := config.ServerConfig.Path + r.RequestURI
	dir := filepath.Dir(fn)
	_ = os.MkdirAll(dir, 0777)
	newFile, err := os.Create(fn)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	_, err = io.Copy(newFile, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}
