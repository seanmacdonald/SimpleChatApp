package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

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
	write_chan := make(chan string)

	go readMessage(read_chan, c)
	go input(write_chan)

	for {
		select {
		case incomingMsg := <-read_chan:
			fmt.Println("Server: ", incomingMsg)
		case outgoingMsg := <-write_chan:
			if err := c.WriteMessage(1, []byte(outgoingMsg)); err != nil {
				fmt.Println(err)
				return
			}
			if outgoingMsg == "bye" {
				return
			}
		}
	}
}

func input(input chan string) {
	defer close(input)

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "bye" {
			input <- text
			return
		}
		input <- text
	}
}

func readMessage(read chan string, c *websocket.Conn) {
	defer close(read)

	for {
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
