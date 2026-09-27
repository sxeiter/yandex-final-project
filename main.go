package main

import (
	"log"
	"os"
	"yandex-final-project/pkg/db"
	"yandex-final-project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {

		log.Fatalf("Критическая ошибка инициализации БД: %v", err)
	}

	defer db.Get().Close()
	server.Run()
}
