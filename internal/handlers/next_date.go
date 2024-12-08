package handlers

import (
	"go_final_project/internal/service"
	"log"
	"net/http"
	"time"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		log.Fatal(err)
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	newDate, err := service.NextDate(now, date, repeat)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(newDate))
	}
}
