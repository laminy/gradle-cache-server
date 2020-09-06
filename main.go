package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/laminy/gradle-cache-server/config"
	"github.com/laminy/gradle-cache-server/handlers"
	"github.com/laminy/gradle-cache-server/middleware"
	"net/http"
)

func main() {
	err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println("Error reading config", err)
		panic(err)
	}
	listenAddr := fmt.Sprintf(":%d", config.ServerConfig.Port)
	router := mux.NewRouter()
	subRouter := router.PathPrefix("").Subrouter()
	subRouter.Methods("PUT").HandlerFunc(handlers.CachePut)
	subRouter.Methods("GET").Handler(http.FileServer(http.Dir(config.ServerConfig.Path)))
	router.Use(middleware.Logging)
	fmt.Printf("Server started at [:%d] with cache at %s\n", config.ServerConfig.Port, config.ServerConfig.Path)
	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
