package main

import (
	"context"

	"github.com/Naman15032001/tolling/types"
)

type GRPCAggregratorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregratorServer(svc Aggregator) *GRPCAggregratorServer {
	return &GRPCAggregratorServer{
		svc: svc,
	}
}

func (s *GRPCAggregratorServer) Aggregate(ctx context.Context, req *types.AggregrateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
