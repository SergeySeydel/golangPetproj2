package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `gorm:"serializer:json"`
	Due  time.Time `json:"due"`
}

type App struct {
	DB *gorm.DB
}

func (a *App) Initialize(dbURI string) {
	db, err := gorm.Open(sqlite.Open(dbURI), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db

	a.DB.AutoMigrate(&Task{})
}

func (a *App) getAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	a.DB.Find(&tasks)
	tasksJSON, _ := json.Marshal(tasks)

	w.WriteHeader(200)
	w.Write([]byte(tasksJSON))
}

func (a *App) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	vars := mux.Vars(r)

	// Select the task with the given id, and convert to JSON.
	result := a.DB.First(&task, "id = ?", vars["id"])
	if result.RowsAffected == 0 {
		http.Error(w, "error: id not found in DataBase", http.StatusNotFound)
		return
	}
	taskJSON, _ := json.Marshal(task)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(taskJSON))
}

func (a *App) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling task create at %s\n", r.URL.Path)
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	rows := a.DB.Create(&newTask).RowsAffected
	log.Println("Added rows: ", rows)

	// Создаем json для ответа
	type ResponseId struct {
		Id int `json:"id"`
	}

	taskJSON, err := json.Marshal(ResponseId{Id: newTask.Id})
	if err != nil {
		http.Error(w, "error: not create task", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(taskJSON))
}

func (a *App) getTaskHandlerByTag(w http.ResponseWriter, r *http.Request) {
	var task Task
	vars := mux.Vars(r)

	
	err := a.DB.Where(&task, "tags = ?", vars["tag"])
	if err.Error != nil {
		http.Error(w, "error: tag not found in DataBase", http.StatusNotFound)
		return
	}
	taskJSON, _ := json.Marshal(task)

	w.WriteHeader(200)
	w.Write([]byte(taskJSON))
}

func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := a.DB.Where("id = ?", vars["id"]).Delete(Task{})
	if err.Error != nil {
		http.Error(w, "error: id not found in DataBase", http.StatusNotFound)

		return
	}

	
	w.WriteHeader(204)
}

func (a *App) deleteAllTaskHandler(w http.ResponseWriter, r *http.Request) {

	err := a.DB.Exec("DELETE FROM tasks")
	if err.Error != nil {
		http.Error(w, "error: not create task", http.StatusInternalServerError)
		return
	}

	
	w.WriteHeader(204)
}

func (a *App) getTaskHandlerByDueDate(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	vars := mux.Vars(r)
	year := vars["yy"]
	month := vars["mm"]
	day := vars["dd"]

	dueDateStr := fmt.Sprintf("%s-%s-%s", year, month, day)
	dueDateTime, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		http.Error(w, "error: invalid date format", http.StatusBadRequest)
		return
	}

	a.DB.Find(&tasks, "due = ?", dueDateTime)
	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "error: failed to serialize tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksJSON)
}

func main() {
	a := &App{}
	a.Initialize("tasks.db")

	r := mux.NewRouter()

	r.HandleFunc("/task", a.getAllTaskHandler).Methods("GET")
	r.HandleFunc("/task/{id}", a.getTaskHandler).Methods("GET")
	r.HandleFunc("/task", a.createTaskHandler).Methods("POST")
	r.HandleFunc("/task/{id}", a.deleteTaskHandler).Methods("DELETE")
	r.HandleFunc("/task", a.deleteAllTaskHandler).Methods("DELETE")
	r.HandleFunc("/due/{yy}/{mm}/{dd}", a.getTaskHandlerByDueDate).Methods("GET")

	r.HandleFunc("/task/tag/{tag}", a.getTaskHandlerByTag).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Listening on 127.0.0.1:1234")
	if err := http.ListenAndServe("127.0.0.1:1234", nil); err != nil {
		panic(err)
	}
}
