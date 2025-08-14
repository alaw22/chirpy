package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"database/sql"
	"sync/atomic"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/alaw22/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func main(){

	const (
		rootPath = "."
		port = "8080"
	)

	// Load environment variables
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}


	dbConn, err := sql.Open("postgres",dbURL)
	if err != nil {
		log.Fatal("Couldn't establish connection to chirpy db: %w",err)
	}


	// Create http handler
	serveMux := http.NewServeMux()
	
	// Create apiconfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: database.New(dbConn),
	}

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
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}




}
