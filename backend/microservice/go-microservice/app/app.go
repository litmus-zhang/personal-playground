package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/greet", GreeterHandler)
	router.HandleFunc("/customers", GetCustomerHandler)
	router.HandleFunc("/customers/{customer_id}", GetCustomer)

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}
