package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go_final_project/internal/models"
)

func (h *Handler) ShowTasks(w http.ResponseWriter, r *http.Request) {
	var errAnswer models.ErrAnswer
	var tasksAnswer models.TasksAnswer
	var searchByDate = false //флаг поиска по дате

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	search := r.FormValue("search")

	// если нет параметра поиска
	if search == "" {
		tasks, err := h.repo.SelectAllTasks()
		tasksAnswer.Tasks = []models.Task{}
		tasksAnswer.Tasks = append(tasksAnswer.Tasks, tasks...)
		if err != nil {
			errAnswer.Error = err.Error()
			resp, err := json.Marshal(errAnswer)
			if err != nil {
				//500 error
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
	} else if date, err := time.Parse("02.01.2006", search); err == nil { // если ищут по дате
		searchByDate = true
		search = date.Format(DateFormat)
		tasks, err := h.repo.SearchTasks(search, searchByDate)
		tasksAnswer.Tasks = []models.Task{}
		tasksAnswer.Tasks = append(tasksAnswer.Tasks, tasks...)
		if err != nil {
			errAnswer.Error = err.Error()
			resp, err := json.Marshal(errAnswer)
			if err != nil {
				//500 error
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
	} else { // если ищут по заголовку или комментарию
		task, err := h.repo.SearchTasks(search, searchByDate)
		tasksAnswer.Tasks = []models.Task{}
		tasksAnswer.Tasks = append(tasksAnswer.Tasks, task...)
		if err != nil {
			errAnswer.Error = err.Error()
			resp, err := json.Marshal(errAnswer)
			if err != nil {
				//500 error
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
	}

	resp, err := json.Marshal(tasksAnswer)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
