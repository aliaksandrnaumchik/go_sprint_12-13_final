package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-scheduler/pkg/db"
	"todo-scheduler/tests"
)

func main() {

	initDb()
	initServer()
}

func getPort() int {
	return tests.Port
}

func initDb() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()
}

func initServer() {
	port := getPort()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Запуск веб‑сервера на порту %d...\n", port)
	fmt.Printf("Откройте в браузере: http://localhost:%d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
