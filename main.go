package main

import (
	// "fmt"
	"net/http"
	"log"
)

func main(){

	// Create http handler
	serveMux := http.NewServeMux()


	// Create server
	server := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	// Start server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start http server")
	}

}
