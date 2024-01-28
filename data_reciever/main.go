package main

import (
	"fmt"
	"github.com/Naman15032001/tolling/types"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

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
	prod  DataProducer
}

func NewDataReciever() (*DataReciever, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
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
	return dr.prod.ProduceData(data)
}
