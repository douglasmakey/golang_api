package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/douglasmakey/backend_base/config"
	m "github.com/douglasmakey/backend_base/middlewares"
	"github.com/douglasmakey/backend_base/databases"
	"github.com/douglasmakey/backend_base/routes"
)

func main() {
	configPath := flag.String("config", "./config/production.json", "path of the config file")
	flag.Parse()

	// Read config
	config, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Init ECHO
	e := echo.New()

	//Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(m.DBMiddleware(databases.Init(config)))

	// Route
	routes.Init(e)

	// Set level logger
	e.Logger.SetLevel(log.INFO)

	// Start server with GraceShutdown
	go func() {
		if err := e.Start(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)); err != nil {
			e.Logger.Info("shutting down the server.")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
