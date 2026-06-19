package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/erfan-goodarzi/booking-event-api/internal/api"
	"github.com/erfan-goodarzi/booking-event-api/internal/routes"
)

// @title           Booking Event API
// @version         1.0
// @description     API for booking events
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	app, err := api.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	routes := routes.RegisterRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("starting on port %d", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
