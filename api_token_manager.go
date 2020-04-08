package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalln("Failed to load .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	GetClient()

	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
