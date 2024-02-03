package client

import (
	"context"
	"github.com/Naman15032001/tolling/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregrateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
