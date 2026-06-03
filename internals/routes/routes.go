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
	protectedRoute.POST("/events", app.Handlers.Event.CreateEvents)
	protectedRoute.PUT("/events/:id", app.Handlers.Event.UpdateEvent)
	protectedRoute.DELETE("/events/:id", app.Handlers.Event.DeleteEvent)
	r.GET("/events", app.Handlers.Event.GetEvents)
	r.GET("/events/:id", app.Handlers.Event.GetEvent)

	// Auth
	r.POST("/signup", app.Handlers.User.Signup)
	r.POST("/login", app.Handlers.User.Login)

	r.GET("/health", app.HealthCheck)

	return r
}
