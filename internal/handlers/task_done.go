package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go_final_project/internal/models"
	"go_final_project/internal/service"
)

func (h *Handler) TaskDone(w http.ResponseWriter, r *http.Request) {
	//var task models.Task
	var errAnswer models.ErrAnswer
	//var buf bytes.Buffer

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, errToInt := strconv.Atoi(r.FormValue("id"))
	task, errToFind := h.repo.SelectTaskById(id)
	if errToInt != nil || errToFind != nil {
		errAnswer.Error = "Неизвестный идентификатор"
		resp, err := json.Marshal(errAnswer)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(resp)
		return
	}

	if task.Repeat == "" {
		err := h.repo.DeleteTask(id)
		if err != nil {
			errAnswer.Error = "Невозможно удалить задачу"
			resp, err := json.Marshal(errAnswer)
			if err != nil {
				//500 error
				log.Println(err.Error())
				http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
	} else {
		newDate, err := service.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errAnswer.Error = err.Error()
			resp, err := json.Marshal(errAnswer)
			if err != nil {
				//500 error
				log.Println(err.Error())
				http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
		task.Date = newDate
		//обновляем данные в таблице
		err = h.repo.UpdateTask(task, id)
		if err != nil {
			errAnswer.Error = err.Error()
			resp, err := json.Marshal(errAnswer)
			if err != nil {
				//500 error
				log.Println(err.Error())
				http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}
	}

	//возвращаем пустой JSON
	resp, err := json.Marshal(struct{}{})
	if err != nil {
		//500 error
		http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
		return
	}
	//ok 200
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
