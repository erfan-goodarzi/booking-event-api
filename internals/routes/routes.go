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
	protectedRoute.POST("/events", app.Handlers.Event.CreateEvents)
	protectedRoute.PUT("/events/:id", app.Handlers.Event.UpdateEvent)
	protectedRoute.DELETE("/events/:id", app.Handlers.Event.DeleteEvent)
	r.GET("/events", app.Handlers.Event.GetEvents)
	r.GET("/events/:id", app.Handlers.Event.GetEvent)

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
