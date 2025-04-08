package models

import (
	"github.com/jinzhu/gorm"
	"github.com/tuanpnt17/eleven-golang-projects/go-bookstore/pkg/config"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:""json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

func (book *Book) CreateBook() *Book {
	db.NewRecord(book)
	db.Create(book)
	return book
}

func DeleteBook(Id int64) *Book {
	var book Book
	db.Where("id = ?", Id).Delete(&book)
	return &book
}

func GetAllBooks() []Book {
	var books []Book
	db.Find(&books)
	return books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var book Book
	db.Where("id = ?", Id).First(&book)
	return &book, db
}
