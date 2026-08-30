package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)


func main() {
	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("GET /static/", http.StripPrefix("/static", fileServer))


	http.HandleFunc(
		"GET /{$}",
		func(writer http.ResponseWriter, request *http.Request) {
			indexTemplate := template.Must(template.ParseFiles("ui/html/index.html"))
			indexTemplate.Execute(writer, nil)
	})

	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
