package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go_final_project/internal/models"
)

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	var errAnswer models.ErrAnswer

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//получаем идетификатор задачи
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		errAnswer.Error = "Неизвестный идентификатор"
		resp, err := json.Marshal(errAnswer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(resp)
		return
	}

	//выполняем запрос к БД
	taskById, err := h.repo.SelectTaskById(id)
	if err != nil {
		log.Println(err.Error())
		errAnswer.Error = "Неизвестный идентификатор"
		resp, err := json.Marshal(errAnswer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(resp)
		return
	}

	resp, err := json.Marshal(taskById)
	if err != nil {
		//500 error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//ok 200
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
