package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_final_project/internal/models"
	"go_final_project/internal/service"
)

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var idAnswer models.IdAnswer
	var errAnswer models.ErrAnswer
	var buf bytes.Buffer

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//читаем данные из тела запроса
	if _, err := buf.ReadFrom(r.Body); err != nil {
		//404 error
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//десериализуем JSON в Task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		//404 error
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

	//проверяем корректность данных из формы
	err := CheckForm(&task)
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

	//добавляем задачу в таблицу и получаем идентификатор записи
	id, err := h.repo.AddTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		//500 error
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//cериализуем полученное id
	idAnswer.ID = id
	resp, err := json.Marshal(idAnswer)
	if err != nil {
		//500 error
		log.Println(err.Error())
		http.Error(w, "Ошибка при сериализации JSON", http.StatusInternalServerError)
		return
	}
	//ok 200
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// функция проверки полей формы и изменения значения date, где это необходимо
func CheckForm(task *models.Task) error {
	if task.Title == "" {
		err := fmt.Errorf("заголовок задачи не может быть пуст")
		return err
	}
	//если дата не указана
	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	}
	//если есть правило повторения
	if task.Repeat != "" {
		newDate, err := service.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
		parsedDate, _ := time.Parse(DateFormat, task.Date) // пропускаем ошибку, тк она вернется в предыдущем шаге
		if parsedDate.Compare(time.Now()) == -1 && task.Date != time.Now().Format(DateFormat) {
			task.Date = newDate
		}
	} else { //если правила нет
		parsedDate, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return err
		}
		if parsedDate.Compare(time.Now()) == -1 {
			task.Date = time.Now().Format(DateFormat)
		}
	}
	return nil
}
