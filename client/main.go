package main 

import (
	"fmt"
	"net/url"
	"log"

	"github.com/gorilla/websocket"
)

func connect(u url.URL) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	
	if err != nil {
		log.Println(err)
	}
	defer c.Close()
	

	for{
		//TODO: handle messages back and forth 
	}
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "connect"}

	fmt.Println("Attempting to connect to: ", u.String())

	connect(u)
}