package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/erfan-goodarzi/booking-event-api/internals/app"
	"github.com/erfan-goodarzi/booking-event-api/internals/routes"
	"github.com/joho/godotenv"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app, err := app.NewApplication()
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
