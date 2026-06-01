package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/erfan-goodarzi/booking-event-api/internals/api"
	"github.com/erfan-goodarzi/booking-event-api/internals/db"
	"github.com/erfan-goodarzi/booking-event-api/internals/models"
	"github.com/erfan-goodarzi/booking-event-api/migrations"
	"github.com/gin-gonic/gin"
)

type Application struct {
	Logger       *log.Logger
	EventHandler *api.EventHandler
	UserHandler  *api.UserHandler
	DB           *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := db.Open()
	if err != nil {
		return nil, err
	}

	err = db.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	eventStore := models.NewPostgresEventStore(pgDB)
	userStore := models.NewPostgresUserStore(pgDB)

	eventHandler := api.NewEventHandler(eventStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)

	app := &Application{
		Logger:       logger,
		EventHandler: eventHandler,
		UserHandler:  userHandler,
		DB:           pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(c *gin.Context) {
	if err := a.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
