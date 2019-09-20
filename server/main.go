package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"runtime"
	"strings"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connect(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	log.Println("A new connection was made. There are now", runtime.NumGoroutine(), "goroutines running.")

	//Note that each http handler func starts a new goroutine and we want to limit
	//this chat application between 2 users: the client and the server.
	//Futhermore, when a websocket connection is terminated the goroutine will terminate.
	if runtime.NumGoroutine() > 3 {
		log.Println("Only one connection allowed at a time")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//call read/write loop
	message(conn)
}

func message(conn *websocket.Conn) {

	read_chan := make(chan string)
	input_chan := make(chan string)
	defer close(read_chan)
	defer close(input_chan)

	go input(input_chan)
	go readMessage(read_chan, conn)

	for {
		select {
		case incomingMsg := <-read_chan:
			fmt.Println(incomingMsg)
		case outgoingMsg := <-input_chan:
			if err := conn.WriteMessage(1, []byte(outgoingMsg)); err != nil {
				fmt.Println(err)
				return
			}

		}
	}
}

func input(c chan string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		c <- strings.TrimSpace(text)
	}
}

func readMessage(c chan string, conn *websocket.Conn) {
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connection to client is over...")
			log.Println(err)
			return
		} else {
			if msgType == 1 {
				c <- "Client: " + string(p)
			}
		}
	}
}

func main() {

	http.HandleFunc("/connect", connect)

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
