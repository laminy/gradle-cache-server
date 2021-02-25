package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/laminy/gradle-cache-server/cache"
	"github.com/laminy/gradle-cache-server/config"
	"github.com/laminy/gradle-cache-server/handlers"
	"github.com/laminy/gradle-cache-server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	err := config.ReadConfig("config.json")
	if err != nil {
		log.Println("Error reading config", err)
		panic(err)
	}
	listenAddr := fmt.Sprintf(":%d", config.ServerConfig.Port)
	router := mux.NewRouter()
	metricsRouter := router.PathPrefix("/metrics").Subrouter()
	metricsRouter.Methods("GET").Handler(promhttp.Handler())
	subRouter := router.PathPrefix("").Subrouter()
	subRouter.Methods("PUT").HandlerFunc(handlers.CachePut)
	subRouter.Methods("GET").Handler(http.FileServer(http.Dir(config.ServerConfig.Path)))
	subRouter.Use(middleware.Logging)
	subRouter.Use(middleware.Metrics)
	subRouter.Use(middleware.FileAccess)
	cache.StartSchedule()
	log.Printf("Server started at [:%d] with cache at %s\n", config.ServerConfig.Port, config.ServerConfig.Path)
	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		panic(err)
	}
}
