package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/Naman15032001/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	listenAddr := ":3000"
	grpclistenAddr := ":3001"
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregrator(store)
	)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(grpclistenAddr, svc)
	makeHTTPTransport(listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP Transport running on port: ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregator(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC Transport running on port: ", listenAddr)
	// make a tcp listener
	ln, err := net.Listen("TCP", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	//make a new grpc server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register our GRPC Server implementation to the GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregratorServer(svc))
	return server.Serve(ln)
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
