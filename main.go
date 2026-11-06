package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)


type Database struct {
	Instance *sql.DB
}
var database = Database{ nil }


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
	DROP TABLE IF EXISTS groceryStores;
	CREATE TABLE IF NOT EXISTS groceryStores (
		id   INTEGER     PRIMARY KEY,
		name VARCHAR(64) NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `
	DROP TABLE IF EXISTS groceryItems;
	CREATE TABLE IF NOT EXISTS groceryItems (
		id       INTEGER     PRIMARY KEY,
		quantity INTEGER     DEFAULT 0,
		name     VARCHAR(64) NOT NULL,
		store    INTEGER     NOT NULL,

		FOREIGN KEY (store) REFERENCES groceryStore(id)
	);`

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

	_, err := db.Exec(query, store.Name)
	if err != nil {
		log.Fatal(err)
	}
}


func addGroceryItem(db *sql.DB, item GroceryItem) {
	query := `
	INSERT INTO groceryItems
	( quantity, name, store )
	VALUES
	( 10, ?, 100 )
	`

	_, err := db.Exec(query, item.Name)
	if err != nil {
		log.Fatal(err)
	}
}


func printGroceryStoresTable(db *sql.DB) {
	query := `
	SELECT * FROM groceryStores
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-20s", "")
	for _, val := range cols {
		fmt.Printf("%-20s", val)
	}
	fmt.Println()

	var id   string
	var name string
	index := 1
	for rows.Next() {
		fmt.Printf("%-20s", strconv.Itoa(index))

		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s%-20s\n", id, name)

		index += 1
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}


func main() {
	fmt.Print("Creating file server...")
	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	fmt.Println("Done!")

	fmt.Print("Opening database...")
	db, err := sql.Open("sqlite3", "./groceries.db")
	if err != nil {
		log.Fatal(err)
	}
	database.Instance = db
	defer db.Close()
	defer func() { database.Instance = nil }()
	fmt.Println("Done!")

	fmt.Println("Testing database...")
	fmt.Println("Creating database...")
	createTables(db)



	printGroceryStoresTable(database.Instance)

	addGroceryStore(db, GroceryStore{ "TestStore", []GroceryItem{} })
	addGroceryStore(db, GroceryStore{ "Meep Mart", []GroceryItem{} })
	addGroceryStore(db, GroceryStore{ "Blue", []GroceryItem{} })
	addGroceryItem(db, GroceryItem{ "TestItem", 10 })

	printGroceryStoresTable(database.Instance)



	http.HandleFunc( "GET /{$}"          , indexHandler  )
	http.HandleFunc( "POST /add-item/{$}", addItemHandler)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func indexHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string][]GroceryStore {
		"GroceryStores": {
			GroceryStore{ "Costco", []GroceryItem{GroceryItem{"Apple"  , 12 }, GroceryItem{ "Banana", 6 }} },
			GroceryStore{ "T&T"   , []GroceryItem{GroceryItem{ "Cereal",  1 }, GroceryItem{ "Bread" , 2 }} },
		},
	}

	tmpl := template.Must(template.ParseFiles(
		"ui/html/index.html",
		"ui/html/grocery-store.html",
	))
	tmpl.Execute(writer, data)
}


func addItemHandler(writer http.ResponseWriter, request *http.Request) {
	itemName  := request.PostFormValue("item-name")
	storeName := request.PostFormValue("store-name")
	log.Println("Requested to add item:")
	log.Println(itemName)
	log.Println(storeName)
}
