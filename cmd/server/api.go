package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Request - struct representing the structure of the HTTP request
type Request struct {
	ID      string            `json:"id"`
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// structure of the response being sent
type Response struct {
	ID      string            `json:"id"`
	status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// Keeping track of the tunnel connections
var tunnelConnections map[string]*websocket.Conn

// mutex to ensure that goroutines dont interrupt
var mutex sync.Mutex

// mapping of pending request to the channels
var responseConnections = make(map[string]chan Response)
var responseMutex sync.Mutex

func handleConnection(c *gin.Context) {
	// we first upgrade to websocket
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

		//  we parse what we recieve as message and send to the waiting connection
		var response Response
		if err := json.Unmarshal(msg, &response); err != nil {
			log.Println("Invalid Repose", err)
			continue
		}

		responseMutex.Lock()
		if ch, exists := responseConnections[response.ID]; exists {
			ch <- response
			delete(responseConnections, response.ID)
		}
		responseMutex.Unlock()

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
	//
	reqID := uuid.New().String()
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		headers[k] = strings.Join(v, ",")
	}
	body := ""

	req := Request{
		ID:      reqID,
		Headers: headers,
		Method:  c.Request.Method,
		Body:    body,
		Path:    c.Request.URL.Path,
	}
	reqData, _ := json.Marshal(req)

	if err = conn.WriteMessage(websocket.TextMessage, reqData); err != nil {
		c.JSON(502, gin.H{"error": "Failed to forward request"})
		return
	}

	// we are waiting for reponse
	respChan := make(chan Response, 1)
	responseMutex.Lock()
	responseConnections[reqID] = respChan
	responseMutex.Unlock()

	select {
	case resp := <-respChan:

		for k, v := range resp.Headers {
			c.Header(k, v)
		}
		c.Data(resp.status, "application/octet-stream", []byte(resp.Body))
	case <-time.After(10 * time.Second):
		c.JSON(504, gin.H{"error": "Request Timeout"})
	}

	// for now, we are working with a dummy response
	// c.JSON(200, gin.H{"message": fmt.Sprintf("Message forwarded to tunnel %s", tunnelID)})
}
