package payment

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// to clean later if not used
type Payment struct {
	Id        int
	ProductId int
	PricePaid float64 // remplasser par une quantité ?
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DEBUT Functions CRUD pour le Payment
func CreatePayment(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

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
		_, errProduct := strconv.ParseUint(productId, 10, 64)

		// check if pricePaid is not a float
		_, errPrice := strconv.ParseFloat(pricePaid, 64)

		if errProduct != nil || errPrice != nil {
			http.Error(w, "productId("+productId+") and pricePaid("+pricePaid+") must be numbers (int/float)", http.StatusBadRequest)
			return
		}

		// Check if price is negative (NEED A REVIEW BC C DEGUEU ALED)
		pricePaidString := fmt.Sprintf("%s", pricePaid) // parse pricePaid to string
		// get first char of pricePaidString
		firstChar := pricePaidString[0:1]
		// check if first char is a minus
		if firstChar == "-" {
			http.Error(w, "pricePaid("+pricePaid+") must be positive", http.StatusBadRequest)
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
			fmt.Fprintf(w, "Payment Created !\n=====================\nPayment ID: %d\nProduct ID: %s\nPrice Paid: %s", id, productId, pricePaid)
		}
	}
}

func UpdatePayment(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

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

func DeletePayment(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

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

func GetPaymentById(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

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

func GetAllPayments(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {

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
