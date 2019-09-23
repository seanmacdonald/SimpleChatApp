package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
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
	if runtime.NumGoroutine() > 4 {
		log.Println("Only one connection allowed at a time")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//start communicating between server and client
	message(conn)
}

func message(conn *websocket.Conn) {

	read_chan := make(chan string)
	input_chan := make(chan string)
	close_input := make(chan int)
	defer close(close_input)

	go input(input_chan, close_input)
	go readMessage(read_chan, conn)

	for {
		select {
		case incomingMsg, ok := <-read_chan:
			if !ok {
				return
			}
			fmt.Println(incomingMsg)
		case outgoingMsg, ok := <-input_chan:
			if !ok {
				return
			}
			if err := conn.WriteMessage(1, []byte(outgoingMsg)); err != nil {
				fmt.Println(err)
				return
			}
			if outgoingMsg == "bye" {
				conn.Close()
				return
			}

		}
	}
}

//TODO: make this exit with the connection function
func input(input chan string, close_in chan int) {
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

func readMessage(read chan string, conn *websocket.Conn) {
	defer close(read)

	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connection to client is over...", runtime.NumGoroutine())
			log.Println(err)
			return
		} else {
			if msgType == 1 {
				read <- "Client: " + string(p)
			}
		}
	}
}

func main() {

	http.HandleFunc("/connect", connect)

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
