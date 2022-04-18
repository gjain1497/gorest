package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func intitializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users", UpdateUser).Methods("PUT")
	r.HandleFunc("/users", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9020", r))
}

func main() {
	InitialMigration()
	intitializeRouter()
}
