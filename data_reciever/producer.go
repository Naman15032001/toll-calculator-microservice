package main

import (
	"encoding/json"
	//"fmt"
	"github.com/Naman15032001/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(types.OBUDATA) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	// start another goroutine to check if we have deliever the data
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					//fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					//fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &KafkaProducer{
		producer: p,
		topic:  topic,
	}, nil
}

func (p *KafkaProducer) ProduceData(data types.OBUDATA) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
}
