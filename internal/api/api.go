package api

import (
	"log"
	"net/http"
	"os"

	"github.com/erfan-goodarzi/booking-event-api/internal/config"
	"github.com/erfan-goodarzi/booking-event-api/internal/db"
	"github.com/erfan-goodarzi/booking-event-api/internal/store"
	"github.com/erfan-goodarzi/booking-event-api/migrations"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Event   *EventHandler
	Ticket  *TicketHandler
	User    *UserHandler
	Booking *BookingHandler
}

type Application struct {
	Logger   *log.Logger
	Handlers *Handler
	DB       db.DB
}

func NewApplication() (*Application, error) {
	cfg := config.Load()
	dsn := cfg.DB.DSN()

	dbConfig := db.DatabaseConfig{Driver: "pgx", DSN: dsn}

	pgDB, err := db.ConnectDB(dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	apiResponse := &APIResponse{}

	eventStore := store.NewPostgresEventStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	ticketStore := store.NewPostgresTicketStore(pgDB)
	bookingStore := store.NewPostgresBookingStore(pgDB)

	eventHandler := NewEventHandler(eventStore, logger, apiResponse)
	userHandler := NewUserHandler(userStore, logger, apiResponse)
	ticketHandler := NewTicketHandler(ticketStore, eventStore, logger, apiResponse)
	bookingHandler := NewBookingHandler(bookingStore, eventStore, logger, apiResponse)

	handlers := &Handler{
		Event:   eventHandler,
		User:    userHandler,
		Ticket:  ticketHandler,
		Booking: bookingHandler,
	}

	app := &Application{
		Logger:   logger,
		Handlers: handlers,
		DB:       pgDB,
	}

	return app, nil
}

// HealthCheck godoc
// @Tags Health
// @Produce json
// @Security BearerAuth
// @Success 200 {object} api.HealthCheckResponse
// @Failure 503 {object} api.HealthCheckErrorResponse
// @Router /health [get]
func (a *Application) HealthCheck(c *gin.Context) {
	if err := a.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
