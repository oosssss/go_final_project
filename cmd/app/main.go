package main

import (
	"fmt"
	"go_final_project/internal/db"
	"go_final_project/internal/handlers"
	"go_final_project/internal/repository"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}

	//получаем порт
	port := os.Getenv("TODO_PORT")
	if port == "" {
		log.Fatal("Необходимо настроить переменную окружения TODO_PORT")
	}

	workDir, _ := os.Getwd()
	//dbFile := filepath.Join(workDir, "scheduler.db")
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		log.Fatal("Необходимо настроить переменную окружения TODO_DBFILE")
	}

	//инициализация БД
	db, err := db.InitDB(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)

	handlerTask := handlers.New(repo)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir(fmt.Sprintf("%s/web/", workDir))))
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Post("/api/task", handlerTask.AddTask)
	r.Get("/api/task", handlerTask.GetTask)
	r.Put("/api/task", handlerTask.EditTask)
	r.Get("/api/tasks", handlerTask.ShowTasks)
	r.Post("/api/task/done", handlerTask.TaskDone)
	r.Delete("/api/task", handlerTask.DeleteTask)

	log.Printf("Сервер запущен на порту %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
