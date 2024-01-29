package main

import (
	"encoding/json"
	"fmt"
	"github.com/Naman15032001/tolling/types"
	"net/http"
)

func main() {
	listenAddr := ":3000"
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregrator(store)
	)
	svc = NewLogMiddleware(svc)
	makeHTTPTransport(listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP Transport running on port: ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregator(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregator(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
