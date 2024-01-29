package main

import (
	"log"
)

// type DistanceCalculator struct{
// 	consumer DataConsumer
// }

const kafkaTopic = "obudata"

// Transport (HTTP ,GRPC , KAFKA) -> attach business logic to consumer to this transport

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	kafkaConsumer, err := NewKafkaconsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
