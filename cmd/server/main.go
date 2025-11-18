package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// we have two main routes
	// The get route for us to connect to the application
	// Any - The core of the application - allow any api request
	r.GET("/connect")
	r.Any("/")

	// Start server on port 8080 (default)
	fmt.Println("...Serving the server")
	r.Run()
}
