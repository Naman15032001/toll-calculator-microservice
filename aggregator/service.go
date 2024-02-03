package main

import (
	"fmt"
	"github.com/Naman15032001/tolling/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregrator struct {
	store Storer
}

func NewInvoiceAggregrator(store Storer) Aggregator {
	return &InvoiceAggregrator{
		store: store,
	}
}

func (i *InvoiceAggregrator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in storage : ", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregrator) CalculateInvoice(obuid int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuid)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuid,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return inv, nil
}


