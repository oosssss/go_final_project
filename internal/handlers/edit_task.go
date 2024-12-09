package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go_final_project/internal/models"
)

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task models.Task
	var errAnswer models.ErrAnswer

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//читаем данные из тела запроса
	if _, err := buf.ReadFrom(r.Body); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//десериализуем JSON в Task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
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

	//проверяем идентификатор записи
	id, errToInt := strconv.Atoi(task.ID)
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

	//проверяем поля формы
	if err := CheckForm(&task); err != nil {
		log.Println(err)
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

	//обновляем данные в таблице
	err := h.repo.UpdateTask(task, id)
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
