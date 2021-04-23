package main

import (
	"log"
	"os"
	"os/signal"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Inicia um serv e habilita o middleware logger
	app := fiber.New()
	app.Use(logger.New())

	// registra os endpoints
	app.Post("/accounts", HandlerCreateAccount)
	app.Get("/accounts/:accountId", HandlerGetAccount)
	app.Post("/transactions", HandlerCreateTransaction)

	// Cria um canal pra capturar cancelamento e interrupção do sistema.
	// a finalidade é encerrar o sistema de modo sutil.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go gracefulShutdown(c, app)

	// Listen and serve boys, listen and serve
	if err := app.Listen(":8000"); err != nil {
		log.Panic(err)
	}

	log.Println("Shutdown complete. Bye bye!")

}

// Aguarda a leitura do signal de cancelamento, e faz shutdown do server
func gracefulShutdown(ch chan os.Signal, app *fiber.App) {
	<-ch
	log.Println("Recebido comando shutdown...")
	err := app.Shutdown()
	if err != nil {
		log.Fatal("Erro ao finalizar aplicação:", err)
		return
	}
}
