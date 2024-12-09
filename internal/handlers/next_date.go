package handlers

import (
	"log"
	"net/http"
	"time"

	"go_final_project/internal/service"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(DateFormat, r.FormValue("now"))
	if err != nil {
		log.Println(err)
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	newDate, err := service.NextDate(now, date, repeat)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(newDate))
}
