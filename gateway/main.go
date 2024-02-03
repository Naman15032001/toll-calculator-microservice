package main

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"net/http"

	"github.com/Naman15032001/tolling/aggregator/client"
	"github.com/sirupsen/logrus"
)

const listenAddr = ":6000"

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	var (
		client         = client.NewHTTPClient("http://localhost:3000")
		invoiceHandler = NewInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("gateway http server running on port %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	//return fmt.Errorf("boom") test error handler
	inv, err := h.client.GetInvoice(context.Background(), 3434)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}
}
