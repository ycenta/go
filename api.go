package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	PricePaid float64 // remplasser par une quantité ?
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DEBUT BLOC BROADCASTER

type Broadcaster interface {
	Register(chan<- interface{})
	Unregister(chan<- interface{})
	Close() error
	Submit(interface{}) bool
}

type broadcaster struct {
	input   chan interface{}
	reg     chan chan<- interface{}
	unreg   chan chan<- interface{}
	outputs map[chan<- interface{}]bool
}

func (bc *broadcaster) broadcast(p interface{}) {
	for ch := range bc.outputs {
		ch <- p
	}
}

func (bc *broadcaster) run() {
	for {
		select {
		case p := <-bc.input:
			bc.broadcast(p)
		case ch, ok := <-bc.reg:
			if ok {
				bc.outputs[ch] = true
			} else {
				return
			}
		case ch := <-bc.unreg:
			delete(bc.outputs, ch)
		}
	}
}

func NewBroadcaster(bufflen int) Broadcaster {
	b := &broadcaster{
		input:   make(chan interface{}, bufflen),
		reg:     make(chan chan<- interface{}),
		unreg:   make(chan chan<- interface{}),
		outputs: make(map[chan<- interface{}]bool),
	}

	go b.run()

	return b

}

func (bc *broadcaster) Register(listener chan<- interface{}) {
	bc.reg <- listener
}

func (bc *broadcaster) Unregister(listener chan<- interface{}) {
	bc.unreg <- listener
}

func (bc *broadcaster) Close() error {
	close(bc.reg)
	close(bc.unreg)
	return nil
}

func (bc *broadcaster) Submit(p interface{}) bool {
	if bc == nil {
		return false
	}
	select {
	case bc.input <- p:
		return true
	default:
		return false
	}
}

// FIN BLOC BROADCASTER

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

	// On init le broadcaster (go routine directement dans le constructeur)
	bc := NewBroadcaster(10)
	// fmt.Println(bc)

	// Nos routes & endpoints API
	r := mux.NewRouter()
	r.HandleFunc("/", getRoot)
	r.HandleFunc("/hello", getHello)

	r.HandleFunc("/product", createProduct(db))              // ici, on passe la variable db à la fonction createProduct (notre handler)
	r.HandleFunc("/product/{id:[0-9]+}", getProductById(db)) //ect
	r.HandleFunc("/product/{id}/update", updateProduct(db))
	r.HandleFunc("/product/{id}/delete", deleteProduct(db))
	r.HandleFunc("/products", getAllProducts(db))

	r.HandleFunc("/payment", createPayment(db, bc))
	r.HandleFunc("/payment/{id:[0-9]+}", getPaymentById(db))
	r.HandleFunc("/payment/{id}/update", updatePayment(db, bc))
	r.HandleFunc("/payment/{id}/delete", deletePayment(db))
	r.HandleFunc("/payments", getAllPayments(db))

	r.HandleFunc("/payment/stream", Stream(bc))

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
func createPayment(db *sql.DB, bc Broadcaster) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) { // closure pour pouvoir utiliser la variable db

		// data from form
		productId := r.FormValue("productId")
		pricePaid := r.FormValue("pricePaid")

		// check inputs emptyness
		if productId == "" || pricePaid == "" {
			http.Error(w, "Missing productId("+productId+") or pricePaid("+pricePaid+")", http.StatusBadRequest)
			return
		}

		// check if productId is not an int
		_, errProduct := strconv.Atoi(productId)

		// check if pricePaid is not a float
		_, errPrice := strconv.ParseFloat(pricePaid, 64)

		if errProduct != nil || errPrice != nil {
			http.Error(w, "productId("+productId+") and pricePaid("+pricePaid+") must be numbers (int/float)", http.StatusBadRequest)
			return
		}

		// check if POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// check if product exists
		productReq := "SELECT * FROM product WHERE id = $1"
		product, err := db.Query(productReq, productId)
		if err != nil {
			fmt.Print("error")
		}
		// if not, return message and stop
		if !product.Next() {
			http.Error(w, "Product does not exist (ID: "+productId+")", http.StatusBadRequest)
			return
		}

		// create payment request
		req := "INSERT INTO payment (\"productId\", \"pricePaid\", \"createdAt\", \"updatedAt\") VALUES ($1, $2, $3, $4) RETURNING id"
		rows, err := db.Query(req, productId, pricePaid, time.Now(), time.Now())
		if err != nil {
			fmt.Print("Error while creating payment (ID: " + productId + ")")
		}

		defer rows.Close()

		// loop through payment rows
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				fmt.Print("No payment was created (ID: " + productId + ")")
			}

			bc.Submit(id)

			fmt.Fprintf(w, "Payment Created !\n=====================\nPayment ID: %d\nProduct ID: %s\nPrice Paid: %s", id, productId, pricePaid)
		}
	}
}

func updatePayment(db *sql.DB, bc Broadcaster) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// check if PUT method
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get payment id
		vars := mux.Vars(r)
		idPayment := vars["id"]

		// get data from form
		productId := r.FormValue("productId")
		pricePaid := r.FormValue("pricePaid")

		// check if inputs are empty
		if productId == "" || pricePaid == "" {
			http.Error(w, "Missing Product ID or Price Paid", http.StatusBadRequest)
			return
		}

		// update payment request
		req := "UPDATE payment SET \"productId\" = $1, \"pricePaid\" = $2, \"updatedAt\" = $3 WHERE id = $4 RETURNING id"
		rows, err := db.Query(req, productId, pricePaid, time.Now(), idPayment)
		if err != nil {
			fmt.Print("Error updating for payment with ID: ", idPayment)
		}

		defer rows.Close()

		// get payment id
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				fmt.Print("No payment was updated")
			}
			fmt.Fprintf(w, "Payment Updated !\n=====================\nPayment ID: %d\nProduct ID: %s\nPrice Paid: %s", id, productId, pricePaid)
		}
	}
}

func deletePayment(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) { // closure pour avoir la variable db

		// check if DELETE method
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get payment id
		vars := mux.Vars(r)
		idPayment := vars["id"]

		fmt.Println(idPayment)

		// delete payment request
		req := "DELETE FROM payment WHERE id = $1"
		_, err := db.Query(req, idPayment)
		if err != nil {
			fmt.Print("Error deleting payment")
		} else {
			fmt.Fprintf(w, "Payment Deleted !\n=====================\nPayment ID: %s", idPayment)
		}

	}
}

func getPaymentById(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

	// closure pour utiliser la variable db
	return func(w http.ResponseWriter, r *http.Request) {

		// check if GET method
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get payment id
		vars := mux.Vars(r)
		idPayment := vars["id"]

		// get payment request
		req := "SELECT * FROM payment WHERE id = $1"
		rows, err := db.Query(req, idPayment)
		if err != nil {
			fmt.Print("Error getting payment with ID: ", idPayment)
		}

		defer rows.Close()

		// get payment id
		for rows.Next() {
			var id int
			var productId int
			var pricePaid string
			var createdAt time.Time
			var updatedAt time.Time
			err = rows.Scan(&id, &productId, &pricePaid, &createdAt, &updatedAt)
			if err != nil {
				fmt.Fprintf(w, "No payment was found with the ID: %d", idPayment)
			}
			fmt.Fprintf(w, "Payment Found !\n=====================\nPayment ID: %d \nProduct ID: %d \nPrice Paid: %s \nCreated At: %s \nUpdated At: %s", id, productId, pricePaid, createdAt, updatedAt)
			return
		}
		fmt.Fprintf(w, "No payment was found with ID: %s", idPayment)
	}

}

func getAllPayments(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) { // closure pour utiliser db

		if r.Method != http.MethodGet { //check GET
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// query + check error
		getRequest := "SELECT * FROM payment"
		rows, err := db.Query(getRequest)
		if err != nil {
			fmt.Print("Error getting all payments")
		}

		defer rows.Close()

		for rows.Next() {
			var id int
			var productId int
			var pricePaid float64
			var createdAt time.Time
			var updatedAt time.Time

			err = rows.Scan(&id, &productId, &pricePaid, &createdAt, &updatedAt)
			if err != nil {
				fmt.Print("No payment was found in DB")
			}

			// Data display
			fmt.Fprintf(w, "Payment ID: %d / Product ID: %d / Price Paid: %0.2f / Created At: %s / Updated At: %s \n", id, productId, pricePaid, createdAt, updatedAt)
		}
	}

}

func Stream(bc Broadcaster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		listener := make(chan interface{})
		bc.Register(listener)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			for p := range listener {
				fmt.Fprintf(w, "Payment: %v\n\n", p)
				wg.Done()
			}

			// bc.Unregister(listener)
		}()
		wg.Wait()
	}
}
