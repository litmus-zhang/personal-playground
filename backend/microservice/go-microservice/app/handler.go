package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"full_name"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

func GreeterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	customers := []Customer{
		{ID: 1,
			Name:    "John Doe",
			Zipcode: "10001"},
		{
			ID:      2,
			Name:    "Jane Doe",
			Zipcode: "10002",
		},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Fprintf(w, "Customer ID: %v\n", vars["customer_id"])
}
