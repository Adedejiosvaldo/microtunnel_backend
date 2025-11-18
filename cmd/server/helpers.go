package main

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func extractTunnelIDFromURL(c *gin.Context) (string, error) {
	pathParts := strings.Split(strings.Trim(c.Request.URL.Path, "/"), "/")

	if len(pathParts) < 1 || pathParts[0] == "" {
		return "", errors.New("Invalid Tunnel ID")
	}
	tunnelID := pathParts[0]
	return tunnelID, nil
}
