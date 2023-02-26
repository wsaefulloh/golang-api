package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// DB set up
func setupDB() *sql.DB {
	host := "127.0.0.1"
	user := "user_pg"
	password := "kode123"
	database := "local_db_psql"
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, database)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Product struct {
	Id            int       `json:"id"`
	Product_name  string    `json:"product_name"`
	Product_price int       `json:"product_price"`
	Product_stock int       `json:"product_stock"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"update_at"`
}

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Product `json:"data"`
	Message string    `json:"message"`
}

// Main function
func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all movies
	router.HandleFunc("/", GetProduct).Methods("GET")

	// Get all movies
	router.HandleFunc("/product", GetProductID).Methods("GET")

	// Create a movie
	router.HandleFunc("/", CreateProduct).Methods("POST")

	// Update a movie
	router.HandleFunc("/", UpdateProduct).Methods("PUT")

	// Delete a specific movie by the movieID
	router.HandleFunc("/{product_id}", DeleteProduct).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 9000")
	err := http.ListenAndServe(":9000", router)

	if err != nil {
		log.Fatal("Error API")
	}
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println(message)
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all movies

// response and request handlers
func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := setupDB()

	printMessage("Getting product...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM products")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var data []Product
	var products Product

	// Foreach movie
	for rows.Next() {

		err := rows.Scan(&products.Id, &products.Product_name, &products.Product_price, &products.Product_stock, &products.Created_at, &products.Updated_at)

		// check errors
		checkErr(err)

		data = append(data, products)
	}

	response := JsonResponse{Type: "success", Data: data}

	json.NewEncoder(w).Encode(&response)
}

// response and request handlers
func GetProductID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := setupDB()

	vars := r.URL.Query()
	productID := vars["product_id"][0]

	fmt.Println("Getting product with ID ", productID)

	// Get all movies from movies table that don't have movieID = "1"
	query := `SELECT * FROM public.products where id = $1`
	rows, err := db.Query(query, productID)

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var data []Product
	var products Product

	// Foreach movie
	for rows.Next() {

		err := rows.Scan(&products.Id, &products.Product_name, &products.Product_price, &products.Product_stock, &products.Created_at, &products.Updated_at)

		// check errors
		checkErr(err)

		data = append(data, products)
	}

	response := JsonResponse{Type: "success", Data: data}

	json.NewEncoder(w).Encode(&response)
}

// response and request handlers
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := JsonResponse{}

	resultProductPrice, err := strconv.Atoi(r.FormValue("product_price"))

	if err != nil {
		response = JsonResponse{Type: "error", Message: "product_price should int"}
		json.NewEncoder(w).Encode(&response)
		return
	}

	resultProductStock, err := strconv.Atoi(r.FormValue("product_stock"))

	if err != nil {
		response = JsonResponse{Type: "error", Message: "product_stock should int"}
		json.NewEncoder(w).Encode(&response)
		return
	}

	products := Product{
		Product_name:  r.FormValue("product_name"),
		Product_price: resultProductPrice,
		Product_stock: resultProductStock,
		Created_at:    time.Now(),
		Updated_at:    time.Now(),
	}

	if products.Product_name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing product_name parameter."}
		json.NewEncoder(w).Encode(&response)
		return
	} else {
		db := setupDB()

		printMessage("Inserting product into DB")

		// fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

		query := `INSERT INTO public.products(product_name, product_price, product_stock, created_at, updated_at) VALUES($1, $2, $3, $4, $5)`
		_, err := db.Exec(query, products.Product_name, products.Product_price, products.Product_stock, products.Created_at, products.Updated_at)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The product has been inserted successfully!"}
		json.NewEncoder(w).Encode(&response)
		return
	}
}

// response and request handlers
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := JsonResponse{}

	resultProductId, err := strconv.Atoi(r.FormValue("product_id"))

	if err != nil {
		response = JsonResponse{Type: "error", Message: "type of product_id must int"}
		json.NewEncoder(w).Encode(&response)
		return
	}

	resultProductPrice, err := strconv.Atoi(r.FormValue("product_price"))

	if err != nil {
		response = JsonResponse{Type: "error", Message: "type of product_price must int"}
		json.NewEncoder(w).Encode(&response)
		return
	}

	resultProductStock, err := strconv.Atoi(r.FormValue("product_stock"))

	if err != nil {
		response = JsonResponse{Type: "error", Message: "type of product_stock must int"}
		json.NewEncoder(w).Encode(&response)
		return
	}

	products := Product{
		Id:            resultProductId,
		Product_name:  r.FormValue("product_name"),
		Product_price: resultProductPrice,
		Product_stock: resultProductStock,
		Updated_at:    time.Now(),
	}

	if products.Product_name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing product_name parameter."}
		json.NewEncoder(w).Encode(&response)
		return
	} else {
		db := setupDB()

		fmt.Println("Update product with ID: ", products.Id)

		// fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)
		query := `UPDATE public.products SET product_name = $1, product_price = $2, product_stock = $3, updated_at = $4  WHERE id = $5`
		_, err := db.Exec(query, products.Product_name, products.Product_price, products.Product_stock, products.Updated_at, products.Id)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The product has been edited successfully!"}
		json.NewEncoder(w).Encode(&response)
		return
	}
}

// response and request handlers
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	productID := params["product_id"]

	response := JsonResponse{}

	if productID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing product)id parameter."}
		json.NewEncoder(w).Encode(&response)
		return
	} else {
		db := setupDB()

		printMessage("Deleting product from DB")

		query := `DELETE FROM public.products where id = $1`

		_, err := db.Exec(query, productID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The product has been deleted successfully!"}
		json.NewEncoder(w).Encode(&response)
		return
	}
}
