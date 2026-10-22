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
	Quantity int
}


var tableNames = [2]string{
	"groceryStores",
	"groceryItems",
}


func arrayContainsValue(arr []string, val string) bool {
	for i := 0; i < len(arr); i++ {
		if val == arr[i] {
			return true
		}
	}
	return false
}


func createTables(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS groceryStores (
		id   INTEGER     PRIMARY KEY,
		name VARCHAR(64) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS groceryItems (
		id       INTEGER     PRIMARY KEY,
		quantity INTEGER     DEFAULT 0,
		name     VARCHAR(64) NOT NULL,
		store    INTEGER,

		FOREIGN KEY(store) REFERENCES groceryStore(id)
	)`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}


func addGroceryStore(db *sql.DB, store GroceryStore) {
	// First check if store exists.
	query := `
	INSERT INTO groceryStores
	( name )
	VALUES
	( ? )
	`

	res, err := db.Exec(query, store.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}


func addGroceryItem(db *sql.DB, item GroceryItem) {
	query := `
	INSERT INTO groceryItems
	( quantity, name, store )
	VALUES
	( 10, ?, 100 )
	`

	res, err := db.Exec(query, item.Name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}


func readFromTable(db *sql.DB, tableName string) *sql.Rows {
	if !arrayContainsValue(tableNames[0:len(tableNames)], tableName) {
		log.Fatal("ERROR (readFromTable): Invalid table name " + tableName)
	}

	query := `
	SELECT * FROM 
	` + tableName

	rows, err := db.Query(query, tableName)
	if err != nil {
		log.Fatal(err)
	}

	return rows
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
	fmt.Println("Creating database...")
	createTables(db)




	fmt.Println("Reading from groceryStores...")
	rows := readFromTable(db, "groceryStores")
	fmt.Println(rows) // ROWS DON'T WORK LIKE THIS GOTTA DO SOMETHING ELSE
	fmt.Println("Reading from groceryItems...")
	rows = readFromTable(db, "groceryItems")
	fmt.Println(rows)

	addGroceryStore(db, GroceryStore{ "TestStore", []GroceryItem{} })
	addGroceryItem(db, GroceryItem{ "TestItem", 10 })

	fmt.Println("Reading from groceryStores...")
	rows = readFromTable(db, "groceryStores")
	fmt.Println(rows)
	fmt.Println("Reading from groceryItems...")
	rows = readFromTable(db, "groceryItems")
	fmt.Println(rows)




	http.HandleFunc( "GET /{$}", indexHandler)
	http.HandleFunc( "POST /add-item/{$}", indexHandler)

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


func addItemHandler(writer http.ResponseWriter, request *http.Request) {
	itemName := request.PostFormValue("item-name")
	storeName := request.PostFormValue("store-name")
	log.Println(itemName)
	log.Println(storeName)
}
