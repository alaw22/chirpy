package main

import (
	"fmt"
	"log"
	"net/http"
)

func main(){

	const (
		rootPath = "."
		port = "8080"
	)

	// Create http handler
	serveMux := http.NewServeMux()
	
	// Create apiconfig
	apiCfg := apiConfig{}

	// FileServer Handler
	fileServerHandler := http.StripPrefix("/app",http.FileServer(http.Dir(rootPath)))

	// Handle requests for files on server. Mapping "/" to the root "."
	serveMux.Handle("/app/",apiCfg.middlewareMetricsInc(fileServerHandler))

	// Register newly defined handlers
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.serverHitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetServerHitsHandler)

	// Create server
	server := &http.Server{
		Addr: ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Serving files from %s at port %s\n",rootPath,port)

	// Start server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}




}
