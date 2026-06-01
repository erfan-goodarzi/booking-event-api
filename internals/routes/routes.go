package routes

import (
	"github.com/erfan-goodarzi/booking-event-api/internals/app"
	"github.com/erfan-goodarzi/booking-event-api/internals/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()
	protectedRoute := r.Group("/")
	protectedRoute.Use(middlewares.Authenticate)

	// EVENTS
	protectedRoute.POST("/events", app.EventHandler.CreateEvents)
	protectedRoute.PUT("/events/:id", app.EventHandler.UpdateEvent)
	protectedRoute.DELETE("/events/:id", app.EventHandler.DeleteEvent)
	r.GET("/events", app.EventHandler.GetEvents)
	r.GET("/events/:id", app.EventHandler.GetEvent)

	// Auth
	r.POST("/signup", app.UserHandler.Signup)
	r.POST("/login", app.UserHandler.Login)

	r.GET("/health", app.HealthCheck)

	return r
}
