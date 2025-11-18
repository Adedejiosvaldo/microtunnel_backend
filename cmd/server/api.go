package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Keeping track of the tunnel connections
var tunnelConnections map[string]*websocket.Conn

// mutex to ensure that goroutines dont interrupt
var mutex sync.Mutex
