package main

import (
	"encoding/json"
	"fmt"
	"github.com/Naman15032001/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// var kafkaTopic string = "obudata"
var kafkaTopic = "obudata"

func main() {
	recv, err := NewDataReciever()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
	fmt.Println("Data Reciever working fine")
}

type DataReciever struct {
	msgch chan types.OBUDATA
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func NewDataReciever() (*DataReciever, error) {
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
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &DataReciever{
		prod: p,
		//msgch: make(chan types.OBUDATA, 128),
	}, nil
}

func (dr *DataReciever) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsRecieveLoop()
}

func (dr *DataReciever) wsRecieveLoop() {
	fmt.Println("new obu client connected")
	for {
		var data types.OBUDATA
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		//fmt.Printf("recieved obu data [%d] <lat:%.2f , long:%.2f> \n", data.OBUID, data.Lat, data.Long)
		if err := dr.produceData(data); err != nil {
			fmt.Println("Produce error", err)
		}
		//dr.msgch <- data
	}
}

func (dr *DataReciever) produceData(data types.OBUDATA) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
	return err
}
