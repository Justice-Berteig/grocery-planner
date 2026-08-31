package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)


type GroceryStore struct {
	Name      string
	Groceries []GroceryItem
}


type GroceryItem struct {
	Name     string
	quantity int
}


func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS testTable (
		id PRIMARY KEY,
		name VARCHAR(20) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}


func addToTable(db *sql.DB) {
	query := `
	INSERT INTO testTable (
		id,
		name
	)
	VALUES (
		1,
		"Justice"
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}


func readFromTable(db *sql.DB) string {
	var id int
	var name string

	query := `
	SELECT * FROM testTable
	`

	err := db.QueryRow(query).Scan(&id, &name)
	if err != nil {
		log.Fatal(err)
	}

	return name
}


func deleteTable(db *sql.DB) {
	query := `
	DROP TABLE IF EXISTS testTable
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}


func main() {
	fmt.Print("Creating file server...")
	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	fmt.Println("Done!")

	fmt.Print("Opening database...")
	db, err := sql.Open("sqlite3", "./groceries.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done!")

	fmt.Println("Testing database...")
	createTable(db)
	addToTable(db)
	res := readFromTable(db)
	fmt.Println(res)
	deleteTable(db)
	fmt.Println("DONE!")

	http.HandleFunc( "GET /{$}", indexHandler)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func indexHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string][]GroceryStore {
		"GroceryStores": {
			GroceryStore{ "Costco", []GroceryItem{GroceryItem{"Apple", 12 } , GroceryItem{ "Banana", 6 }} },
			GroceryStore{ "T&T"   , []GroceryItem{GroceryItem{ "Cereal", 1 }, GroceryItem{ "Bread", 2 }} },
		},
	}

	tmpl := template.Must(template.ParseFiles(
		"ui/html/index.html",
		"ui/html/grocery-store.html",
	))
	tmpl.Execute(writer, data)
}
