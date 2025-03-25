package main

import (
	"github.com/gorilla/pat"
	"net/http"
)

func main() {
	router := pat.New()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("dist"))).Methods("GET")
	http.ListenAndServe(":8000", router)
}
