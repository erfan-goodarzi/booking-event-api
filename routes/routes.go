package routes

import (
	"github.com/erfan-goodarzi/booking-event-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	protectedRoute := server.Group("/")
	protectedRoute.Use(middlewares.Authenticate)

	// EVENTS
	protectedRoute.POST("/events", createEvents)
	protectedRoute.PUT("/events/:id", updateEvent)
	protectedRoute.DELETE("/events/:id", deleteEvent)
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	// Auth
	server.POST("/signup", signup)
	server.POST("/login", login)
}
