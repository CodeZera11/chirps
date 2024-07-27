package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	const port = "8080"

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	apiCfg := apiConfig{}

	mux.Handle("/app/*",	apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer((http.Dir("."))))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCheckHits)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileServerHits int
}

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerCheckHits(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// resp := fmt.Sprintf("Hits: %v", cfg.fileServerHits)
	html := fmt.Sprintf(`<html>
		<body>
    	<h1>Welcome, Chirpy Admin</h1>
    	<p>Chirpy has been visited %v times!</p>
		</body>
		</html>`, cfg.fileServerHits)

		w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileServerHits = 0
}