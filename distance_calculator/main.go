package main

import (
	"log"

	"github.com/Naman15032001/tolling/aggregator/client"
)

// type DistanceCalculator struct{
// 	consumer DataConsumer
// }



const (
	kafkaTopic = "obudata"
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
	
	kafkaConsumer, err := NewKafkaconsumer(kafkaTopic, svc , client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
