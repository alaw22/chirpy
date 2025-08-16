package main

import (
	"fmt"
	"net/http"
)



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

// func (cfg *apiConfig) middlewareAdminReset(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){

// 	})
// }

func (cfg *apiConfig) serverHitsHandler(w http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,cfg.fileserverHits.Load())

	w.Header().Set("Content-Type","text/html")

	w.WriteHeader(200)
	w.Write([]byte(body))
}

// func (cfg *apiConfig) resetServerHitsHandler(w http.ResponseWriter, req *http.Request) {
// 	cfg.fileserverHits.Store(0)
// 	w.WriteHeader(200)
// 	w.Write([]byte("Successfully reset server hits"))
// }

func (cfg *apiConfig) resetUsersAndHitsHandler(w http.ResponseWriter, req *http.Request) {
	// Only local dev access
	if cfg.platform != "dev"{
		respondWithError(w,
						  403,
						  "Forbidden",
						  nil)
		return
	}

	// Deleting all users from database
	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil {
		respondWithError(w,
						  505,
						  "PostgreSQL Error: Couldn't Delete Users.",
						  err)
		return
	}

	// Reset server hits 
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Successfully reset server hits and deleted all users in database"))

}
