package routes

import (
	"github.com/gorilla/mux"
	"github.com/tuanpnt17/eleven-golang-projects/go-bookstore/pkg/controllers"
	"net/http"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/books/", controllers.CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/books/", controllers.GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/books/{bookId}", controllers.GetBookById).Methods(http.MethodGet)
	router.HandleFunc("/books/{bookId}", controllers.UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/books/{bookId}", controllers.DeleteBook).Methods(http.MethodDelete)
}
