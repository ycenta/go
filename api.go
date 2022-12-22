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

// Constantes de connexion à la base de données
const (
	host     = "localhost"
	port     = 5431
	user     = "postgres"
	password = "example"
	dbname   = "postgres"
)

// Fonctions pour les routes (a ignorer, truc de la doc)
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

// Nos objets
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

	// Pour info, pour pouvoir utiliser l'objet db, il faut le passer dans les handlers, mais les fonctions utilisent des closures parce que c'est plus simple pour acceder à la variable db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Notre connexion à la base de données
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	fmt.Println("Connected to database")
	defer db.Close()

	// Nos routes & endpoints API
	r := mux.NewRouter()
	r.HandleFunc("/", getRoot)
	r.HandleFunc("/hello", getHello)

	r.HandleFunc("/product", createProduct(db))              // ici, on passe la variable db à la fonction createProduct (notre handler)
	r.HandleFunc("/product/{id:[0-9]+}", getProductById(db)) //ect
	r.HandleFunc("/product/{id}/update", updateProduct(db))
	r.HandleFunc("/product/{id}/delete", deleteProduct(db))
	r.HandleFunc("/products", getAllProducts(db))

	r.HandleFunc("/payment", createPayment)
	r.HandleFunc("/payment/{id}", getPaymentById)
	r.HandleFunc("/payment/{id}/update", updatePayment)
	r.HandleFunc("/payment/{id}/delete", deletePayment)
	r.HandleFunc("/payments", getAllPayments)

	// On lance le serveur
	err = http.ListenAndServe(":3333", r)
	if err != nil {
		fmt.Println(err)
	}

}

// DEBUT Functions CRUD pour le Product
func createProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { // ici, closure pour pouvoir utiliser la variable db

		// On récupère les données du formulaire
		name := r.FormValue("name")
		price := r.FormValue("price")

		if r.Method != http.MethodPost { //On check si c'est EN POST
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if name == "" || price == "" {
			http.Error(w, "Missing data", http.StatusBadRequest)
			return
		}

		insertRequest := "INSERT INTO product (name, price, \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4) RETURNING id"
		rows, err := db.Query(insertRequest, name, price, time.Now(), time.Now())
		if err != nil {
			fmt.Print("error creating product")
		}
		defer rows.Close()
		// Comme dans la requete d'insert, on lui dit "retourne moi l'id du produit", on va récuperer les infos dans "rows" et la renvoyer dans la réponse
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				fmt.Print("nothing was created")
			}
			fmt.Fprintf(w, "%d", id)
		}

	}
}

func updateProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { // ici, closure pour pouvoir utiliser la variable db

		if r.Method != http.MethodPut { //On check si c'est EN PUT
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		//Parametres de l'url, et variables du formulaire
		vars := mux.Vars(r)
		id_product := vars["id"]
		name := r.FormValue("name")
		price := r.FormValue("price")

		if name == "" || price == "" {
			http.Error(w, "Missing data", http.StatusBadRequest)
			return
		}

		updateRequest := "UPDATE product SET name = $1, price = $2, \"updatedAt\" = $3 WHERE id = $4 RETURNING id"
		rows, err := db.Query(updateRequest, name, price, time.Now(), id_product)
		if err != nil {
			fmt.Print("error updating product")
		}
		defer rows.Close()
		// Comme dans la requete d'update, on lui dit "retourne moi l'id du produit", on va récuperer les infos dans "rows" et la renvoyer dans la réponse
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				fmt.Print("nothing was updated")
			}
			fmt.Fprintf(w, "%d", id)
		}
	}
}

func deleteProduct(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { // ici, closure pour pouvoir utiliser la variable db
		if r.Method != http.MethodDelete { //On check si c'est EN DELETE
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		id_product := vars["id"]

		deleteRequest := "DELETE FROM product WHERE id = $1"
		_, err := db.Query(deleteRequest, id_product)
		if err != nil {
			fmt.Print("error")
		} else {
			fmt.Fprintf(w, "success")
		}

	}
}

func getProductById(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { // ici, closure pour pouvoir utiliser la variable db

		if r.Method != http.MethodGet { //On check si c'est EN GET
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		id_product := vars["id"]

		getRequest := "SELECT * FROM product WHERE id = $1"
		rows, err := db.Query(getRequest, id_product)
		if err != nil {
			fmt.Print("error")
		}
		defer rows.Close()
		// Comme dans la requete d'insert, on lui dit "retourne moi l'id du produit", on va récuperer les infos dans "rows" et la renvoyer dans la réponse*
		// if rows is empty, it will return nothing
		for rows.Next() {
			var id int
			var name string
			var price string
			var createdAt time.Time
			var updatedAt time.Time
			err = rows.Scan(&id, &name, &price, &createdAt, &updatedAt)
			if err != nil {
				fmt.Print("nothing was found")
			}
			fmt.Fprintf(w, "%d %s %s %s %s", id, name, price, createdAt, updatedAt)
			return
		}
		fmt.Fprintf(w, "nothing was found")

	}
}

func getAllProducts(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { // ici, closure pour pouvoir utiliser la variable db

		if r.Method != http.MethodGet { //On check si c'est EN GET
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		getRequest := "SELECT * FROM product"
		rows, err := db.Query(getRequest)
		if err != nil {
			fmt.Print("error")
		}
		defer rows.Close()
		// Comme dans la requete d'insert, on lui dit "retourne moi l'id du produit", on va récuperer les infos dans "rows" et la renvoyer dans la réponse*
		// if rows is empty, it will return nothing
		for rows.Next() {
			var id int
			var name string
			var price string
			var createdAt time.Time
			var updatedAt time.Time
			err = rows.Scan(&id, &name, &price, &createdAt, &updatedAt)
			if err != nil {
				fmt.Print("nothing was foundzz")
			}
			fmt.Fprintf(w, "%d %s %s %s %s", id, name, price, createdAt, updatedAt)
		}

	}
}

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
