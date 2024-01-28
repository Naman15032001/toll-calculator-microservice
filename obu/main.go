package main

import (
	"fmt"
	"github.com/Naman15032001/tolling/types"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"time"
)

const wsEndPoint = "ws://localhost:30000/ws"

var sendInterval = time.Second

func genLatLong() (float64, float64) {
	return genCords(), genCords()
}

func genCords() float64 {
	n := (float64)(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	obuIDS := genOBUIDS(20)
	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUDATA{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%+v\n", data)
		}
		time.Sleep(sendInterval)
	}

}

func genOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < len(ids); i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
