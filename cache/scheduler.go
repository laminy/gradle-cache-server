package cache

import (
	"github.com/laminy/gradle-cache-server/config"
	"github.com/laminy/gradle-cache-server/handlers"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func StartSchedule() {
	log.Printf("Cache cleanup will be started every %s and cache alive %s\n", config.ServerConfig.Scan, config.ServerConfig.Alive)
	ticker := time.NewTicker(config.ServerConfig.ScanInterval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				scan()
				break
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	scan()
}

func scan() {
	log.Println("Start scanning cache for delete")
	err := filepath.Walk(config.ServerConfig.Path, verifyFile)
	if err != nil {
		log.Println(err)
	}
	log.Println("Cache scan complete")
}

func verifyFile(path string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return err
	}
	fileAccessTime := time.Unix(info.Sys().(*syscall.Stat_t).Atim.Unix())
	diff := time.Now().Sub(fileAccessTime)
	if config.ServerConfig.AliveInterval > diff {
		return nil
	}
	relPath := strings.Replace(path, config.ServerConfig.Path, "", 1)
	log.Printf("Deleting file %s\n", relPath)
	os.Remove(path)
	handlers.IncCounter(relPath, 200, "DELETE")
	return nil
}
