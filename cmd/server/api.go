package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Keeping track of the tunnel connections
var tunnelConnections map[string]*websocket.Conn

// mutex to ensure that goroutines dont interrupt
var mutex sync.Mutex

func handleConnection(c *gin.Context) {
	// we first upgrade
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println("websocket upgrade")
		c.JSON(500, gin.H{"error": "websocket upgrade failed"})
	}

	defer conn.Close()

	tunnelID := uuid.New().String()[:15]

	// mutex lock, save the id, and unlock
	mutex.Lock()
	tunnelConnections[tunnelID] = conn
	mutex.Unlock()

	// send back the tunnel info to the cli
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Tunnel established: %s", tunnelID)))

	// This is where we listen for connection and then respond if we have anything
	// also by keeping the server open
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading websocket", err)
			// then we delete the connection, we lock it to ensure another routine doesnt work on it
			// while we are at it
			mutex.Lock()
			delete(tunnelConnections, tunnelID)
			mutex.Unlock()
			break
		}
		log.Printf("Recieved message from CLI Client: %s: %s", tunnelID, msg)

	}

}
