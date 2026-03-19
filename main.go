package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-scheduler/pkg/api"
	"todo-scheduler/pkg/db"
	"todo-scheduler/tests"
)

func main() {
	initDb()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}()

	initApi()
	initServer()
}

func initDb() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
}

func initApi() {
	api.Init()
}

func initServer() {
	port := tests.Port
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Запуск веб‑сервера на порту %d...\n", port)
	fmt.Printf("Откройте в браузере: http://localhost:%d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
