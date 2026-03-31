package main

import (
	"flag"

	"github.com/bartek5186/waitlistfox/controllers"
	"github.com/bartek5186/waitlistfox/internal"
	"github.com/bartek5186/waitlistfox/internal/i18n"
	mid "github.com/bartek5186/waitlistfox/internal/middleware"
	"github.com/bartek5186/waitlistfox/services"
	"github.com/bartek5186/waitlistfox/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var migrate bool

func init() {
	flag.BoolVar(&migrate, "migrate", false, "Run DB migrations on start")
}

func main() {
	flag.Parse()

	logger := internal.NewLogger()

	if err := i18n.LoadTranslations(); err != nil {
		logrus.WithError(err).Fatal("failed to load translations")
	}

	config := internal.LoadConfiguration("config/config.json")
	db := internal.NewDatabaseConnection(config)

	if migrate {
		if err := internal.MigrateRun(db); err != nil {
			logrus.WithError(err).Fatal("failed to run migrations")
		}
	}

	appStore := store.NewAppStore(db, &config)
	appService := services.NewAppService(appStore, logger.GetLogger())
	waitlistController := controllers.NewWaitlistController(appService, logger.GetLogger())

	e := echo.New()
	e.HideBanner = true
	e.Validator = internal.NewInputValidator()
	e.Use(mid.LanguageMiddleware("pl", config.Languages))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.Domains,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))

	e.Static("/", "static")
	e.GET("/health", waitlistController.Health)
	e.POST("/v1/waitlist/subscribe", waitlistController.Subscribe)
	e.POST("/v1/waitlist/unsubscribe", waitlistController.Unsubscribe)

	e.Logger.Fatal(e.Start(config.ServerAddress()))
}
