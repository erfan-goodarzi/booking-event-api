package routes

import (
	_ "github.com/erfan-goodarzi/booking-event-api/docs"
	"github.com/erfan-goodarzi/booking-event-api/internals/app"
	"github.com/erfan-goodarzi/booking-event-api/internals/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()
	protectedRoute := r.Group("/")
	protectedRoute.Use(middlewares.Authenticate)

	// EVENTS
	events := r.Group("/events")
	{
		events.GET("", app.Handlers.Event.GetEvents)
		events.GET("/:id", app.Handlers.Event.GetEvent)
	}

	protectedEvents := protectedRoute.Group("/events")
	{
		protectedEvents.POST("", app.Handlers.Event.CreateEvent)
		protectedEvents.POST("/:id", app.Handlers.Event.DeleteEvent)
		protectedEvents.PUT("/:id", app.Handlers.Event.UpdateEvent)
		protectedEvents.POST("/:id/tickets", app.Handlers.Ticket.CreateTicket)
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
