package main

import (
	"github.com/erfan-goodarzi/booking-event-api/db"
	"github.com/erfan-goodarzi/booking-event-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
