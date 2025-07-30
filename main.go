package main

import (
	"fmt"
	"net/http"
	"log"
)

func main(){

	const (
		rootPath = "."
		port = "8080"
	)

	// Create http handler
	serveMux := http.NewServeMux()
	
	// Handle root requests
	serveMux.Handle("/",http.FileServer(http.Dir(rootPath)))

	// Create server
	server := &http.Server{
		Addr: ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Serving files from %s at port %s\n",rootPath,port)

	// Start server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start http server")
	}




}
