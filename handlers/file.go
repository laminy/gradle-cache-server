package handlers

import (
	"github.com/laminy/gradle-cache-server/config"
	"os"
	"time"
)

func ModifyFileAccessTime(path string, code int, method string) {
	if method != "GET" && code != 200 {
		return
	}
	fn := config.ServerConfig.Path + path
	info, err := os.Stat(fn)
	if err != nil {
		return
	}
	now := time.Now().Local()
	info.Sys()
	os.Chtimes(fn, now, info.ModTime())
}
