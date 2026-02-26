package main

import (
	"fmt"
	"net/http"
)

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
