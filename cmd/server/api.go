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

// Handles every request from the client: cli
func handleRequest(c *gin.Context) {

	tunnelID, err := extractTunnelIDFromURL(c)

	if tunnelID == "" || err != nil {
		c.JSON(400, gin.H{"error": "Invalid Tunnel ID"})
	}

	// check if we already have that connection
	//  and add mutex
	mutex.Lock()
	conn, exists := tunnelConnections[tunnelID]
	mutex.Unlock()

	if !exists {
		c.JSON(400, gin.H{"error": "Tunnel not found or has been disconnected"})
		return
	}
	// Serialize the request
	requestData := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)

	err = conn.WriteMessage(websocket.TextMessage, []byte(requestData))

	if err != nil {
		c.JSON(502, gin.H{"error": "Failed to forward request"})
		return
	}

	// for now, we are working with a dummy response
	c.JSON(200, gin.H{"message": fmt.Sprintf("Message forwarded to tunnel %s", tunnelID)})
}
