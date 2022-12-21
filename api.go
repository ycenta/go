package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5431
	user     = "postgres"
	password = "example"
	dbname   = "postgres"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

type Product struct {
	Id        int
	Name      string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Payment struct {
	Id        int
	ProductId int
	PricePaid float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {

	// Pour info, pour pouvoir utiliser l'objet db, il faut le passer dans les handlers, mais les fonctions utilisent des closures
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	fmt.Println("Connected to database")
	defer db.Close()

	r := mux.NewRouter()
	// Nos routes
	r.HandleFunc("/", getRoot)
	r.HandleFunc("/hello", getHello)
	r.HandleFunc("/product", createProduct(db))
	r.HandleFunc("/product/{id:[0-9]+}", getProductById)
	r.HandleFunc("/product/{id}/update", updateProduct)
	r.HandleFunc("/product/{id}/delete", deleteProduct)
	r.HandleFunc("/products", getAllProducts)

	r.HandleFunc("/payment", createPayment)
	r.HandleFunc("/payment/{id}", getPaymentById)
	r.HandleFunc("/payment/{id}/update", updatePayment)
	r.HandleFunc("/payment/{id}/delete", deletePayment)
	r.HandleFunc("/payments", getAllPayments)

	err = http.ListenAndServe(":3333", r)
	if err != nil {
		fmt.Println(err)
	}

}

// DEBUT Functions CRUD pour le Product
func createProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { // ici, closure pour pouvoir utiliser la variable db

		if r.Method != http.MethodPost { //On check si c'est EN POST
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		insertRequest := "INSERT INTO product (name, price, \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4) RETURNING id"
		rows, err := db.Query(insertRequest, "test", 1.99, time.Now(), time.Now())
		if err != nil {
			fmt.Print("error creating product")
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				fmt.Print("nothing was created")
			}
			// return id to webpage
			fmt.Fprintf(w, "%d", id)
		}

		// fmt.Fprintf(w, "Produit créé ! (on renverra le produzit)")
	}
}

func updateProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut { //On check si c'est EN PUT
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Produit modifié ! (on renverra le produit)")
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete { //On check si c'est EN DELETE
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Produit supprimé !")
}

func getProductById(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet { //On check si c'est EN GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the db variable from the main func
	// row := db.QueryRow("SELECT * FROM product WHERE id = ?", 1)
	// stringifyrow := fmt.Sprintf("%s", row)
	fmt.Fprintf(w, "Produit récupéré ! (on renverra le produit)")
	// fmt.Fprintf(w, stringifyrow)
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet { //On check si c'est EN GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Produits récupérés ! (on renverra les produits)")
}

// FIN Functions CRUD pour le Product

// DEBUT Functions CRUD pour le Payment
func createPayment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost { //On check si c'est EN POST
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Payment created")
}

func updatePayment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut { //On check si c'est EN PUT
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Payment updated")
}

func deletePayment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete { //On check si c'est EN DELETE
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Payment deleted")
}

func getPaymentById(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet { //On check si c'est EN GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Payment retrieved")
}

func getAllPayments(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet { //On check si c'est EN GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Payment retrieved")
}
