package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/form", handleForm)

	fmt.Printf("Starting server on :8686\n")
	if err := http.ListenAndServe(":8686", nil); err != nil {
		log.Fatal(err)
	}
}

func handleForm(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		_, _ = fmt.Fprintf(writer, "ParseForm() error: %v", err)
	}
	name := request.FormValue("name")
	address := request.FormValue("address")
	_, _ = fmt.Fprintf(writer, "Name: %s\nAddress: %s\n", name, address)
}

func handleHello(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/hello" {
		http.Error(writer, "404 not found", http.StatusNotFound)
		return
	}
	if request.Method != http.MethodGet {
		http.Error(writer, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := fmt.Fprintf(writer, "Hello world")
	if err != nil {
		return
	}
}
