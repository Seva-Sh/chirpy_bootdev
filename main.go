package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Seva-Sh/chirpy_bootdev/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// stateful config
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

// middleware pattern that wraps handlers and increments the counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	port := "8080"
	mux := http.NewServeMux()

	// get db_url from the environment
	dbURL := os.Getenv("DB_URL")

	// get PLATFORM evironment
	plt := os.Getenv("PLATFORM")

	// open connection to a database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open the Database Connection:", err)
	}

	dbQueries := database.New(db)

	// initialize api config
	apiCfg := &apiConfig{db: dbQueries, platform: plt}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// ensure the program loggs an error if encountered with server listening failure
	log.Fatal(srv.ListenAndServe())
}
