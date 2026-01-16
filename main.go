package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	payload := map[string]string{"message": "Hello, World"}

	mux.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}))

	log.Fatal(http.ListenAndServe(":3001", mux))
}
