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
	

	message(c) 
}

func message(c *websocket.Conn) {
	read_chan := make(chan string)

	go readMessage(read_chan, c)

	for {
		select {
		case incomingMsg := <-read_chan:
			fmt.Println(incomingMsg)
		}
	}
}

func readMessage(read chan string, c *websocket.Conn) {
	defer close(read)

	for{
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		read <- string(message)
	}
	
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "connect"}

	fmt.Println("Attempting to connect to: ", u.String())

	connect(u)
}