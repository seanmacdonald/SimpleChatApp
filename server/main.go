package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"runtime"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connect(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	fmt.Println("A new connection was made. There are now", runtime.NumGoroutine(), "goroutines running.")
	
	//Note that each http handler func starts a new goroutine and we want to limit
	//this chat application between 2 users: the client and the server
	if runtime.NumGoroutine() > 3 {
		fmt.Println("Only one connection allowed at a time")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("message type: ", msgType)
			if msgType == 1 {
				fmt.Println("message: ", string(p))
			}

			//echo the message back now
			if err := conn.WriteMessage(1, p); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func main() {

	http.HandleFunc("/connect", connect)

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
