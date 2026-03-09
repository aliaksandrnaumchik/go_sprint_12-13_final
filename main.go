package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-scheduler/tests"
)

func main() {
	port := tests.Port

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Запуск веб‑сервера на порту %d...\n", port)
	fmt.Printf("Откройте в браузере: http://localhost:%d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
