package main

import (
	"fmt"
	"net/http"
	"log"
)

func readinessHandler(w http.ResponseWriter, req *http.Request){
	
	body := "OK"

	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.WriteHeader(200)	
	w.Write([]byte(body))

}

func main(){

	const (
		rootPath = "."
		port = "8080"
	)

	// Create http handler
	serveMux := http.NewServeMux()
	
	// Handle requests for files on server. Mapping "/" to the root "."
	serveMux.Handle("/app/",http.StripPrefix("/app",http.FileServer(http.Dir(rootPath))))

	// Register a handler
	serveMux.HandleFunc("/healthz", readinessHandler)

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
