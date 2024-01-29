package main

import (
	"fmt"
	"github.com/Naman15032001/tolling/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
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
