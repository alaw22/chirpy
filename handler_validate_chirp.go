package main

import (
	"io"
	"encoding/json"
	"net/http"
	"strings"
)

func replace_profanity(msg string) string{
	const replacement = "****"
	var profane_words []string = []string{"kerfuffle","sharbert","fornax"}

	profane_word_indexes := make([]int,0)
	split_msg := strings.Split(msg," ")

	for i, word := range split_msg{
		for _, profane_word := range profane_words{
			if strings.ToLower(word) == profane_word{
				profane_word_indexes = append(profane_word_indexes,i)
				break
			}
		}
	}
	
	for _, idx := range profane_word_indexes{
		split_msg[idx] = replacement
	}

	return strings.Join(split_msg," ")

}

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	
	type chirpBody struct {
		Body string `json:"body"`
	}
	
	// type jsonValid struct {
	// 	Valid bool `json:"valid"`
	// }
	type cleanChirp struct {
		Clean_Body string `json:"cleaned_body"`
	}
	

	dat, err := io.ReadAll(req.Body)
	if err != nil {
		respondeWithError(w,
						  501,
						  "Couldn't read request body",
						  err)
		return
	}


	chirp := chirpBody{}


	err = json.Unmarshal(dat,&chirp)
	if err != nil {

		respondeWithError(w,
						  502,
						  "Couldn't unpack json to chirpBody{}",
						  err)
						  
		return
	}

	if len(chirp.Body) > 140 {
		respondeWithError(w,
						  400,
						  "Chirp is too long",
						  nil)
	} else {

		// Clean chirp body
		cleanChirp := cleanChirp{
			Clean_Body: replace_profanity(chirp.Body),
		}

		respondWithJSON(w, 200, cleanChirp)
		
	}

}