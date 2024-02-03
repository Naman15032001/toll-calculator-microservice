package main

import "github.com/Naman15032001/tolling/types"

type GRPCAggregratorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregratorServer(svc Aggregator) *GRPCAggregratorServer {
	return &GRPCAggregratorServer{
		svc : svc,
	}
}

func (s *GRPCAggregratorServer) AggregateDistance(req *types.AggregrateRequest) error {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix : req.Unix,
	}
	return s.svc.AggregateDistance(distance)
}