package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/bt-order-svc/models"
)

type VerifyJWT func(token string) AuthResponse

type OrderServer struct {
	store     models.OrderStore
	verifyJWT VerifyJWT
}

func (o *OrderServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header["Token"] == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResponse := o.verifyJWT(r.Header["Token"][0])
	if authResponse.Status == INVALID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if authResponse.Status == NOT_FOUND {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/order/all" {
			o.getAllOrders(w, r, authResponse)
		} else if r.URL.Path == "/order/current" {
			o.getCurrentOrders(w, r, authResponse)
		}
	}
}

func (o *OrderServer) getAllOrders(w http.ResponseWriter, r *http.Request, authResponse AuthResponse) {
	orders, _ := o.store.GetOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(orders)
}

func (o *OrderServer) getCurrentOrders(w http.ResponseWriter, r *http.Request, authResponse AuthResponse) {
	orders, _ := o.store.GetCurrentOrdersByCustomerID(authResponse.ID)
	json.NewEncoder(w).Encode(orders)
}
