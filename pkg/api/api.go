package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	log.Printf("Received request: now=%s, date=%s, repeat=%s", nowStr, dateStr, repeat)

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "Invalid now parameter", http.StatusBadRequest)
		return
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		log.Printf("Error processing request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "") // возвращаем пустую строку при ошибке
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, result)
}
