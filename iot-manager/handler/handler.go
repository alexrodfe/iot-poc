// Package handler exposes an HTTP REST API for managing IoT devices.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the HTTP routes for the IoT manager.
func RegisterRoutes(router *gin.Engine) {
	router.GET("/devices/:id", getDeviceEntriesByID)
}

func getDeviceEntriesByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id})
}
