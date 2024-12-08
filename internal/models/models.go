package models

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type IdAnswer struct {
	ID string `json:"id"`
}

type ErrAnswer struct {
	Error string `json:"error"`
}

type TasksAnswer struct {
	Tasks []Task `json:"tasks"`
}
