package main

import (
	"log"
	"os"
	"os/signal"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lrweck/go-sampleapi/http"
	"github.com/lrweck/go-sampleapi/repository"
	pg "github.com/lrweck/go-sampleapi/repository/postgresql"
	"github.com/lrweck/go-sampleapi/service"
)

func main() {

	pgRepo, err := pg.NewPGRepo("postgres://postgres:postgres@127.0.0.1:5432/sampleapi?application_name=api", 30)
	if err != nil {
		log.Fatalf("%s", err)
	}
	repo := repository.NewApiRepository(pgRepo)
	apiServ := service.NewApiService(repo)
	handler := http.NewHandler(apiServ)

	// Inicia um serv e habilita o middleware logger
	app := fiber.New()
	app.Use(logger.New())

	// registra os endpoints
	app.Post("/accounts", handler.PostAccount)
	app.Get("/accounts/:accountId", handler.GetAccount)
	app.Post("/transactions", handler.PostTransaction)

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
