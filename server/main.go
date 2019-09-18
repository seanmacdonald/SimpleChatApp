package main 

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024, 
	WriteBufferSize: 1024, 
}

func connect(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err :=  upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}	

	for {
		msgType, p, err :=  conn.ReadMessage() 
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("message type: ", msgType)
			if msgType == 1 {
				fmt.Println("message: ", string(p))
			}
		}
	}
}

func main() { 

	http.HandleFunc("/connect", connect)	

	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)

}
