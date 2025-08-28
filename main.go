package main

import (
	"fmt"
	"os"
	"log"
	"time"
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
	platform string
	secretKey string
	polkaAPIKey string
}

const defaultExpirationTime = time.Second*3600 // seconds -> hour

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
	
	secretKey := os.Getenv("SECRET_STRING")
	if secretKey == ""{
		log.Fatal("Secret key is empty cannot provide JWTs without key")
	}

	polkaAPIKey := os.Getenv("POLKA_KEY")
	if polkaAPIKey == ""{
		log.Fatal("Polka api key is empty. Cannot verify transactions")
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
		platform: os.Getenv("PLATFORM"),
		secretKey: secretKey,
		polkaAPIKey: polkaAPIKey,
	}

	// FileServer Handler
	fileServerHandler := http.StripPrefix("/app",http.FileServer(http.Dir(rootPath)))

	// Handle requests for files on server. Mapping "/" to the root "."
	serveMux.Handle("/app/",apiCfg.middlewareMetricsInc(fileServerHandler))

	// Register newly defined handlers
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.serverHitsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetUsersAndHitsHandler)
	serveMux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	serveMux.HandleFunc("PUT /api/users",apiCfg.updateUserHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpHandler)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpHandler)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshHandler)
	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUserHandler)

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
