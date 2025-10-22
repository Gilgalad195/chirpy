package main

import (
	"fmt"
	"net/http"
)

func (c *apiConfig) countHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl :=
		`<html>
		  <html>
		    <h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`
	html := fmt.Sprintf(tmpl, c.fileserverHits.Load())
	w.Write([]byte(html))
}
