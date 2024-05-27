/*
// Пример REST сервера с несколькими маршрутами(используем только стандартную библиотеку)

// POST   /task/              :  создаёт задачу и возвращает её ID
// GET    /task/<taskid>      :  возвращает одну задачу по её ID
// GET    /task/              :  возвращает все задачи
// DELETE /task/<taskid>      :  удаляет задачу по ID
// DELETE /task/              :  удаляет все задачи
// GET    /tag/<tagname>      :  возвращает список задач с заданным тегом
// GET    /due/<yy>/<mm>/<dd> :  возвращает список задач, запланированных на указанную дату

*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"proj/internal/taskstore"

	"github.com/gorilla/mux"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	log.Printf("NewTaskServer Constructor started\n")
	return &taskServer{store: store}
}

func trimIDFromRequest(r *http.Request) (int, error) {
	log.Printf("trimIDFromRequest started\n")

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		//http.Error(w, "expect 'task/<id>' in task handler", http.StatusBadRequest)
		return 0, fmt.Errorf("expect 'task/<id>' in task handler")
	}
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return 0, err
	}
	log.Printf("trimIDFromRequest returned id: %d\n", id)
	return id, nil
}

func trimDateFromRequest(r *http.Request) (int, time.Month, int, error) {
	vars := mux.Vars(r)

	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid year parameter")
	}

	monthInt, err := strconv.Atoi(vars["month"])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid month parameter")
	}
	if monthInt < 1 || monthInt > 12 {
		return 0, 0, 0, fmt.Errorf("invalid month parameter")
	}
	month := time.Month(monthInt)

	day, err := strconv.Atoi(vars["day"])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid day parameter")
	}

	return year, month, day, nil
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling task create at %s\n", r.URL.Path)

	// Структура для создания задачи
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	// Для ответа в виде JSON
	type ResponseId struct {
		Id int `json:"id"`
	}

	// JSON в качестве Content-Type
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	// Обработка тела запроса
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Создаем новую задачу в хранилище и получаем ее <id>
	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)

	// Создаем json для ответа
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (ts *taskServer) getAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get all tasks at %s\n", r.URL.Path)

	allTasks := ts.store.GetAllTasks()

	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get task at %s\n", r.URL.Path)

	id, err := trimIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // код ошибки 404
		return
	}
	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling delete task at %s\n", r.URL.Path)

	id, err := trimIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	erro := ts.store.DeleteTask(id)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusNotFound) // код ошибки 404
		return
	}
}

func (ts *taskServer) deleteAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling delete all tasks at %s\n", r.URL.Path)

	ts.store.DeleteAllTasks()
}

func (ts *taskServer) getTaskByDue(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling delete all tasks at %s\n", r.URL.Path)

	year, month, day, err := trimDateFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение задач по указанной дате
	allTasks := ts.store.GetTasksByDueDate(year, month, day)

	// Преобразование задач в JSON
	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTasksByTag(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get tasks by tag at %s\n", r.URL.Path)

	vars := mux.Vars(r)
	tag := vars["tagname"]

	// Получение задач по указанному тегу
	tagTasks := ts.store.GetTasksByTag(tag)

	// Преобразование задач в JSON
	js, err := json.Marshal(tagTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	r := mux.NewRouter()
	server := NewTaskServer()

	// Added routing for "/task/" URL
	r.HandleFunc("/task/", server.createTaskHandler).Methods("POST")
	r.HandleFunc("/task/", server.getAllTaskHandler).Methods("GET")
	r.HandleFunc("/task/", server.deleteAllTaskHandler).Methods("DELETE")
	r.HandleFunc("/task/{id:[0-9]+}", server.getTaskHandler).Methods("GET")
	r.HandleFunc("/task/{id:[0-9]+}", server.deleteTaskHandler).Methods("DELETE")
	r.HandleFunc("/task/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", server.getTaskByDue).Methods("GET")
	r.HandleFunc("/tag/{tagname}", server.getTasksByTag).Methods("GET")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
