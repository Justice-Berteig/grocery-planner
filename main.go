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
		storeID  INTEGER     NOT NULL,

		FOREIGN KEY (storeID) REFERENCES groceryStore(id)
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
	( quantity, name, storeID )
	VALUES
	( ?, ?, 100 )
	`

	_, err := db.Exec(query, item.Quantity, item.Name)
	if err != nil {
		log.Fatal(err)
	}
}


func printGroceryStoresTable(db *sql.DB) {
	query := `
	SELECT
		store.id,
		store.name,
		COUNT(*) OVER (PARTITION BY item.id) AS Items
	FROM
		   groceryStores store INNER JOIN groceryItems item
		ON store.id = item.storeID
	`

	// Run query and get rows
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print column names
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-20s", "")
	for _, val := range cols {
		fmt.Printf("%-20s", val)
	}
	fmt.Println()

	// Scan and print row values
	var id        string
	var name      string
	var itemCount string
	index := 1
	for rows.Next() {
		fmt.Printf("%-20s", strconv.Itoa(index))

		if err := rows.Scan(&id, &name, &itemCount); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s%-20s%-20s\n", id, name, itemCount)

		index += 1
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
}


func printGroceryItemsTable(db *sql.DB) {
	query := `
	SELECT
		item.id,
		item.name,
		item.quantity,
		store.name AS 'Store Name'
	FROM
		   groceryItems item INNER JOIN groceryStores store
		ON item.storeID = store.id
	`

	// Run query and get rows
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print column names
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-20s", "")
	for _, val := range cols {
		fmt.Printf("%-20s", val)
	}
	fmt.Println()

	// Scan and print row values
	var id        string
	var name      string
	var quantity  string
	var storeName string

	index := 1
	for rows.Next() {
		fmt.Printf("%-20s", strconv.Itoa(index))

		if err := rows.Scan(&id, &name, &quantity, &storeName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s%-20s%-20s%-20s\n", id, name, quantity, storeName)

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
	fmt.Println()
	printGroceryItemsTable(database.Instance)

	addGroceryStore(db, GroceryStore{ "TestStore", []GroceryItem{} })
	addGroceryStore(db, GroceryStore{ "Meep Mart", []GroceryItem{} })
	addGroceryStore(db, GroceryStore{ "Blue", []GroceryItem{} })
	addGroceryItem(db, GroceryItem{ "Apple", 5 })
	addGroceryItem(db, GroceryItem{ "Orange", 10 })
	addGroceryItem(db, GroceryItem{ "Cereal", 2 })
	addGroceryItem(db, GroceryItem{ "Bread", 20 })
	addGroceryItem(db, GroceryItem{ "Crackers", 12 })
	addGroceryItem(db, GroceryItem{ "Chicken", 6 })

	printGroceryStoresTable(database.Instance)
	fmt.Println()
	printGroceryItemsTable(database.Instance)



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
