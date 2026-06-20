package routes

import (
	_ "github.com/erfan-goodarzi/booking-event-api/docs"
	"github.com/erfan-goodarzi/booking-event-api/internal/api"
	"github.com/erfan-goodarzi/booking-event-api/internal/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(app *api.Application) *gin.Engine {
	r := gin.Default()
	protectedRoute := r.Group("/")
	protectedRoute.Use(middlewares.Authenticate)

	// EVENTS
	events := r.Group("/events")
	{
		events.GET("", app.Handlers.Event.GetAll)
		events.GET("/:id", app.Handlers.Event.GetById)
	}

	protectedEvents := protectedRoute.Group("/events")
	{
		protectedEvents.POST("", app.Handlers.Event.Create)
		protectedEvents.POST("/:id", app.Handlers.Event.Delete)
		protectedEvents.PUT("/:id", app.Handlers.Event.Update)
		protectedEvents.POST("/:id/tickets", app.Handlers.Ticket.Create)
		protectedEvents.POST("/tickets/:id/register", app.Handlers.Booking.RegisterEvent)
		protectedEvents.PUT("/tickets/register/:id/status", app.Handlers.Booking.UpdateRegistrationStatus)
	}

	// Playlist
	playlist := protectedRoute.Group("/playlist")
	{
		playlist.POST("", app.Handlers.Playlist.Create)
		playlist.GET("", app.Handlers.Playlist.GetAll)
		playlist.GET("/:id", app.Handlers.Playlist.GetById)
		playlist.PUT("/:id", app.Handlers.Playlist.Update)
		playlist.DELETE("/:id", app.Handlers.Playlist.Delete)
		playlist.POST("/:playlistId/events/:eventId", app.Handlers.Playlist.AddEvent)
	}

	// Auth
	auth := r.Group("/auth")
	{
		auth.POST("/signup", app.Handlers.User.Signup)
		auth.POST("/login", app.Handlers.User.Login)
		auth.POST("/refresh", app.Handlers.User.Refresh)
		auth.POST("/logout", app.Handlers.User.Logout)
	}

	r.GET("/health", app.HealthCheck)

	//doc
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
