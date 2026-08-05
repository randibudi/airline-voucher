package main

import (
	"log"

	"github.com/randibudi/airline-voucher/backend/internal/app"
)

func main() {
	application := app.New()
	log.Fatal(application.Listen(":8080"))
}
