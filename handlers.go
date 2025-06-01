package main

import (
	"encoding/json"
	"fmt"
	"goserver/internal/database"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQuerries     *database.Queries
}

type parameters struct {
	Body string `json:"body"`
}

type cleanedVals struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorVals struct {
	Error string `json:"error"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respBody := errorVals{
			Error: "Something went wrong",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(500)
		w.Write(dat)
	}

	if len(params.Body) > 140 {
		log.Println("Chirp is too long")
		respBody := errorVals{
			Error: "Chirp is too long",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	respBody := clean_body(params.Body)

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	html_body := fmt.Sprintf(`<html> 
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html_body))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func clean_body(body string) cleanedVals {
	profaneWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	var cleaned string

	splitted := strings.Split(body, " ")

	for i, split := range splitted {
		_, exists := profaneWords[strings.ToLower(split)]
		if exists {
			if i == 0 {
				fmt.Println("first")
				cleaned = "****"
				continue
			}
			cleaned = cleaned + " ****"
			continue
		}
		if i == 0 {
			fmt.Println("first")
			cleaned = split
			continue
		}
		cleaned = cleaned + " " + split
	}

	respBody := cleanedVals{
		CleanedBody: cleaned,
	}

	return respBody
}
