package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	// Nos routes
	r.HandleFunc("/", getRoot)
	r.HandleFunc("/hello", getHello)
	r.HandleFunc("/product", createProduct)
	r.HandleFunc("/product/{id:[0-9]+}", getProductById)
	r.HandleFunc("/product/{id}/update", updateProduct)
	r.HandleFunc("/product/{id}/delete", deleteProduct)
	r.HandleFunc("/products", getAllProducts)

	r.HandleFunc("/payment", createPayment)
	r.HandleFunc("/payment/{id}", getPaymentById)
	r.HandleFunc("/payment/{id}/update", updatePayment)
	r.HandleFunc("/payment/{id}/delete", deletePayment)
	r.HandleFunc("/payments", getAllPayments)

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		fmt.Println(err)
	}

}

// DEBUT Functions CRUD pour le Product
func createProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost { //On check si c'est EN POST
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Produit créé ! (on renverra le produit)")
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

	fmt.Fprintf(w, "Produit récupéré ! (on renverra le produit)")
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
