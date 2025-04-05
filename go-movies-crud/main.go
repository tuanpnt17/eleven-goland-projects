package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

func main() {
	router := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "1234567890", Title: "The Matrix", Director: &Director{1, "Tuan", "Pham"}})
	movies = append(movies, Movie{ID: "2", Isbn: "0987654321", Title: "The Adventure", Director: &Director{2, "Anh", "Vo"}})

	router.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	router.HandleFunc("/movie/{id}", getMovie).Methods(http.MethodGet)
	router.HandleFunc("/movies", createMovie).Methods(http.MethodPost)
	router.HandleFunc("/movie/{id}", updateMovie).Methods(http.MethodPut)
	router.HandleFunc("/movie/{id}", deleteMovie).Methods(http.MethodDelete)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func deleteMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	movieId := params["id"]

	for i, movie := range movies {
		if movie.ID == movieId {
			movies = append(movies[:i], movies[i+1:]...)
			break
		}
	}

	_ = json.NewEncoder(writer).Encode(movies)
}

func updateMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	movieId := params["id"]
	var requestMovie Movie
	_ = json.NewDecoder(request.Body).Decode(&requestMovie)
	for _, movie := range movies {
		if movie.ID == movieId {
			movie.Isbn = requestMovie.Isbn
			movie.Title = requestMovie.Title
			movie.Director = requestMovie.Director
			_ = json.NewEncoder(writer).Encode(movie)
			return
		}
	}
}

func createMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(request.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	_ = json.NewEncoder(writer).Encode(movie)
}

func getMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	movieId := params["id"]
	for _, movie := range movies {
		if movie.ID == movieId {
			_ = json.NewEncoder(writer).Encode(movie)
			return
		}
	}
}

func getMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(movies)
	if err != nil {
		return
	}
}
