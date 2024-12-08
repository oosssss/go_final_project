package handlers

import (
	"encoding/json"
	"go_final_project/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var errAnswer models.ErrAnswer
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, errToInt := strconv.Atoi(r.FormValue("id"))
	_, errToFind := h.repo.SelectTaskById(id)
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
