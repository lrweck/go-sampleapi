package main

import (
	"log"
	"os"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lrweck/go-sampleapi/db"
)

func main() {

	// Starts a new Mux
	app := fiber.New()
	app.Use(logger.New())

	// // Create os.Signal channel to intercept errors and interruptions
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)

	// go graceful(c, app)

	// if err := app.Listen(":8000"); err != nil {
	// 	log.Panic(err)
	// }

	// log.Println("Shutdown complete. Bye bye!")

	conn, err := db.GetConn()

	var opeid int
	var desc string

	rows, err := conn.Query(`SELECT * FROM operationtypes`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&opeid, &desc); err != nil {
			log.Fatal(err)
		}
		log.Printf("ID: %d - Desc: %s\n", opeid, desc)
	}

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
