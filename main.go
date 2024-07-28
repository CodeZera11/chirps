package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codezera11/chirps/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB *database.DB
}

func main() {
	const filerootpath = "."
	const port = "8080"

	db, err := database.NewDB("database.json")

	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		DB: db,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/app/*",	apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer((http.Dir(filerootpath))))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCheckHits)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetOneChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
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