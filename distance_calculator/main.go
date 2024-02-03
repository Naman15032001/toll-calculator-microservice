package main

import (
	"github.com/Naman15032001/tolling/aggregator/client"
	"log"
)

// type DistanceCalculator struct{
// 	consumer DataConsumer
// }

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000/aggregate"
)

// Transport (HTTP ,GRPC , KAFKA) -> attach business logic to consumer to this transport

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	//httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer, err := NewKafkaconsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
