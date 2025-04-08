package main

import (
	"github.com/gorilla/mux"
	"github.com/tuanpnt17/eleven-golang-projects/go-bookstore/pkg/routes"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8686", nil))
}
