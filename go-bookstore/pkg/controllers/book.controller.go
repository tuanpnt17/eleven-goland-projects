package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tuanpnt17/eleven-golang-projects/go-bookstore/pkg/models"
	"github.com/tuanpnt17/eleven-golang-projects/go-bookstore/pkg/utils"
	"net/http"
	"strconv"
)

func GetBooks(writer http.ResponseWriter, request *http.Request) {
	newBooks := models.GetAllBooks()
	res, _ := json.Marshal(newBooks)
	writer.Header().Set("Content-Type", "pkglication/json")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write(res)
	if err != nil {
		return
	}
}

func GetBookById(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error occurred when parsing book id")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	newBook, _ := models.GetBookById(ID)
	res, _ := json.Marshal(newBook)
	writer.Header().Set("Content-Type", "pkglication/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		return
	}
}

func CreateBook(writer http.ResponseWriter, request *http.Request) {
	createBook := &models.Book{}
	utils.ParseBody(request, createBook)
	book := createBook.CreateBook()
	res, _ := json.Marshal(book)
	writer.Header().Set("Content-Type", "pkglication/json")
	writer.WriteHeader(http.StatusCreated)
	_, err := writer.Write(res)
	if err != nil {
		return
	}
}

func UpdateBook(writer http.ResponseWriter, request *http.Request) {
	var updatedBook = &models.Book{}
	utils.ParseBody(request, updatedBook)
	vars := mux.Vars(request)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error occurred when parsing book id")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	bookDetails, db := models.GetBookById(ID)
	if updatedBook.Name != "" {
		bookDetails.Name = updatedBook.Name
	}
	if updatedBook.Author != "" {
		bookDetails.Author = updatedBook.Author
	}
	if updatedBook.Publication != "" {
		bookDetails.Publication = updatedBook.Publication
	}
	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	writer.Header().Set("Content-Type", "pkglication/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		return
	}
}

func DeleteBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		fmt.Println("Error occurred when parsing book id")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	deletedBook := models.DeleteBook(ID)
	res, _ := json.Marshal(deletedBook)
	writer.Header().Set("Content-Type", "pkglication/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		return
	}
}
