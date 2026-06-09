package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/erfan-goodarzi/booking-event-api/internals/api"
	"github.com/erfan-goodarzi/booking-event-api/internals/db"
	"github.com/erfan-goodarzi/booking-event-api/internals/store"
	"github.com/erfan-goodarzi/booking-event-api/migrations"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Event  *api.EventHandler
	Ticket *api.TicketHandler
	User   *api.UserHandler
}

type Application struct {
	Logger   *log.Logger
	Handlers *Handler
	DB       *sql.DB
}

func NewApplication() (*Application, error) {
	dbSrc := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	pgDB, err := db.ConnectDB(dbSrc)
	if err != nil {
		return nil, err
	}

	err = db.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	apiResponse := &api.APIResponse{}

	eventStore := store.NewPostgresEventStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	ticketStore := store.NewPostgresTicketStore(pgDB)

	eventHandler := api.NewEventHandler(eventStore, logger, apiResponse)
	userHandler := api.NewUserHandler(userStore, logger, apiResponse)
	ticketHandler := api.NewTicketHandler(ticketStore, eventStore, logger, apiResponse)

	handlers := &Handler{
		Event:  eventHandler,
		User:   userHandler,
		Ticket: ticketHandler,
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
