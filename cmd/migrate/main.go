package main

import (
	"log"

	"subscriptions/internal/infrastructure/migration"
)

func main() {
	if err := migration.Up(nil); err != nil {
		log.Fatal(err)
	}
}
