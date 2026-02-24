package main

import (
	"database/sql"
	"fmt"
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
}

// middleware pattern that wraps handlers and increments the counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// handler that writes the number of request that have been counted
func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
	val32 := cfg.fileserverHits.Load()
	htmlContent := fmt.Sprintf(`
	<html>
	  <body>
	    <h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	  </body>
	</html>`, val32)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlContent))
}

// handler that resets apiConfig count to 0
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func main() {
	godotenv.Load()
	port := "8080"
	mux := http.NewServeMux()

	// get db_url from the environment
	dbURL := os.Getenv("DB_URL")

	// open connection to a database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open the Database Connection:", err)
	}

	dbQueries := database.New(db)

	// initialize api config
	apiCfg := &apiConfig{db: dbQueries}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// ensure the program loggs an error if encountered with server listening failure
	log.Fatal(srv.ListenAndServe())
}
