package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"yandex-final-project/pkg/api"
)

func Run() {

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	api.Init()

	addr := ":" + port
	fmt.Printf("🚀 Сервер запущен! Откройте в браузере: http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
