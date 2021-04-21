package main

import (
	"log"
	"os"
	"os/signal"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Starts a new Mux
	app := fiber.New()
	app.Use(logger.New())

	app.Post("/accounts", HandlerCreateAccount)
	app.Get("/accounts/:accountId", HandlerGetAccount)
	app.Post("/transactions", HandlerCreateTransaction)

	// Create os.Signal channel to intercept errors and interruptions
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go graceful(c, app)

	if err := app.Listen(":8000"); err != nil {
		log.Panic(err)
	}

	log.Println("Shutdown complete. Bye bye!")

}

func graceful(ch chan os.Signal, app *fiber.App) {
	<-ch
	log.Println("Received shutdown command...")
	err := app.Shutdown()
	if err != nil {
		log.Fatal("Error shutting down:", err)
		return
	}
}
