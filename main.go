package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//WebSockets are upgraded HTTP connections that live until the connection is killed by either the client or the server.
// It’s through this WebSocket connection that we can perform duplex communication which is a really fancy way of saying
// we can communicate to-and-from the server from our client using this single connection.

//The real beauty of WebSockets is that they use a grand total of 1 TCP connection and all communication is done over this
//single long-lived TCP connection. This drastically reduces the amount of network overhead required to build real-time
//applications using WebSockets as there isn’t a constant polling of HTTP endpoints required.

//We'll need to define our Upgrader
//this will require a read and write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Creating a simple HTTP server

func main() {

	fmt.Println("Goo WebSockets !!")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	//says i want to allow any connection to my HTTP endpoint regardless
	// of what the origin of that connection is
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// helpful log statement to show connections
	log.Println("Client Connected")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/ws", wsEndpoint)
}

//define a reader which will listen for
//new messages being sent to our websocket endpoint
func reader(conn *websocket.Conn) {
	for {
		//read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		//print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
