package main

import (
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
)

func readinessHandler(w http.ResponseWriter, req *http.Request){
	
	body := "OK"

	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)	
	w.Write([]byte(body))

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) serverHitsHandler(w http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Hits: %v",cfg.fileserverHits.Load())
	w.WriteHeader(200)
	w.Write([]byte(body))
}

func (cfg *apiConfig) resetServerHitsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Successfully reset server hits"))
}


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
	serveMux.HandleFunc("/healthz", readinessHandler)
	serveMux.HandleFunc("/metrics", apiCfg.serverHitsHandler)
	serveMux.HandleFunc("/reset", apiCfg.resetServerHitsHandler)

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
