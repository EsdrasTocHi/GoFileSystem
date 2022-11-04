package api

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewApi() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/holamundo", HolaMundo).Methods("GET")
	r.HandleFunc("/mkdisk", mkdisk).Methods("GET")
	r.HandleFunc("/command", ReadCommand).Methods("POST")
	r.HandleFunc("/kafka", kafka).Methods("POST")
	log.Println("RestAPI up")
	log.Fatal(http.ListenAndServe(":3030", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
