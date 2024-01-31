package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Naman15032001/tolling/types"
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
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "missing OBU ID",
			})
			return
		}
		obuid, err := strconv.Atoi(values[0])
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid OBU ID",
			})
			return
		}
		invoice, err := svc.CalculateInvoice(obuid)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		WriteJSON(w, http.StatusOK, invoice)
	}

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
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
