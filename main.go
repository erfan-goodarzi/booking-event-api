package main

import (
	"net/http"

	"example.com/booking-event/db"
	"example.com/booking-event/models"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	server.GET("/events", getEvents)
	server.POST("/events", createEvents)
	server.Run(":8080")
}

func getEvents(context *gin.Context) {
	events := models.GetAllEvents()
	context.JSON(http.StatusOK, events)
}

func createEvents(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request"})
		return
	}
	event.Save()
	context.JSON(http.StatusCreated, gin.H{"msg": "Success", "event": event})
}
