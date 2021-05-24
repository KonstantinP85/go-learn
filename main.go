package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type Reader struct {
	gorm.Model
	Name     string
	Email    string `gorm:"type varchar(100);unique_index"`
	Age      int
	Books    []Book
}

type Book struct {
	gorm.Model
	Title    string
	Author   string
	Year     int
	ReaderId int
}
var (
	reader = &Reader{Name: "Jon", Email: "jon@test.ru", Age: 23}
	books = []Book{
		{Title: "First book", Author: "author1", Year: 2010, ReaderId: 1},
		{Title: "Second book", Author: "author2", Year: 2015, ReaderId: 1},
	}
)

var db *gorm.DB
var err error

func main() {
	dialect := "postgres"
	host := "localhost"
	dbPort := "5432"
	user := "postgres"
	dbName := "go_project"
	password := "Kristina"

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	db, err = gorm.Open(dialect, dbURI)
	if err !=nil {
		log.Fatal(err)
	} else {
		fmt.Println("successfully")
	}

	defer db.Close()

	db.AutoMigrate(&Reader{})
	db.AutoMigrate(&Book{})

	router := mux.NewRouter()
	router.HandleFunc("/readers", getReaders). Methods("GET")
	router.HandleFunc("/books", getBooks). Methods("GET")
	router.HandleFunc("/reader/{id}", getOneReader). Methods("GET")
	router.HandleFunc("/book/{id}", getOneBook). Methods("GET")
	router.HandleFunc("/reader/create", createReader). Methods("POST")
	router.HandleFunc("/book/create", createBook). Methods("POST")
	router.HandleFunc("/reader/update/{id}", updateReader).Methods("PUT")
	router.HandleFunc("/book/update/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/reader/delete/{id}", deleteReader). Methods("DELETE")
	router.HandleFunc("/book/delete/{id}", deleteBook). Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Reader controllers

func getReaders(w http.ResponseWriter, r *http.Request) {
	var readers []Reader
	db.Find(&readers)
	json.NewEncoder(w).Encode(&readers)
}

func getOneReader(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var oneReader Reader
	var books []Book
	db.First(&oneReader, params["id"])
	db.Model(&oneReader).Related(&books)
	oneReader.Books = books
	json.NewEncoder(w).Encode(&oneReader)
}

func createReader(w http.ResponseWriter, r *http.Request) {
	var reader Reader
	json.NewDecoder(r.Body).Decode(&reader)

	createdReader := db.Create(&reader)
	err = createdReader.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&reader)
	}
}

func updateReader(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var reader Reader
	db.First(&reader, params["id"])
	json.NewDecoder(r.Body).Decode(&reader)
	db.Save(&reader)
	json.NewEncoder(w).Encode(reader)
}

func deleteReader(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var reader Reader
	db.First(&reader, params["id"])
	db.Delete(&reader)
	json.NewEncoder(w).Encode(&reader)
}

//book controllers

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(&books)
}

func getOneBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var oneBook Book
	db.First(&oneBook, params["id"])
	json.NewEncoder(w).Encode(&oneBook)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.First(&book, params["id"])
	json.NewDecoder(r.Body).Decode(&book)
	db.Save(&book)
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.First(&book, params["id"])
	db.Delete(&book)
	json.NewEncoder(w).Encode(&book)
}