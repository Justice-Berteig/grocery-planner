package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)


type GroceryStore struct {
	Name      string
	Groceries []string
}


func main() {
	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	http.HandleFunc( "GET /{$}", indexHandler)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func indexHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string][]GroceryStore {
		"GroceryStores": {
			GroceryStore{ "Costco", []string{"Apple", "Banana"} },
			GroceryStore{ "T&T"   , []string{"Cereal", "Bread"} },
		},
	}

	tmpl := template.Must(template.ParseFiles(
		"ui/html/index.html",
		"ui/html/grocery-store.html",
	))
	tmpl.Execute(writer, data)
}
