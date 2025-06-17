package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	TodoList = []Todo{}
	mu       sync.Mutex
	nexId    = 1
)

func main() {

	http.HandleFunc("/todo", TodoHandler)
	http.HandleFunc("/todo/", TodoByIdHandler)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST, DELETE")
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {

	case http.MethodGet:
		log.Println(Log("GET todo/"))
		json.NewEncoder(w).Encode(TodoList)

	case http.MethodPost:
		var newTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			log.Println(Log("POST todo/ Invalid item format").WithLevel(LevelError))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if newTodo.Item == "" {
			log.Println(Log("POST todo/<todo_id> Item cannot be empty").WithLevel(LevelWarning))
			http.Error(w, "Item cannot be empty", http.StatusBadRequest)
			return
		}

		todo := Todo{ID: nexId, Item: newTodo.Item}
		nexId++
		mu.Lock()
		TodoList = append(TodoList, todo)
		mu.Unlock()
		log.Println(Log("POST todo/"))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	case http.MethodDelete:
		log.Println(Log("DELETE todo/"))
		TodoList = []Todo{}
		json.NewEncoder(w).Encode(map[string]string{"message": "Todo cleared successfully"})
	default:
		log.Println(Log("todo/ Method not allowed").WithLevel(LevelWarning))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func TodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	todoId := strings.TrimPrefix(r.URL.Path, "/todo/")
	id, err := strconv.Atoi(todoId)
	if err != nil {
		log.Println(Log("todo/<todo_id> Invalid todo id").WithLevel(LevelError))
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	var todo *Todo
	for index := range TodoList {
		if TodoList[index].ID == id {
			todo = &TodoList[index]
			break
		}
	}

	if todo == nil {
		log.Println(Log("todo/<todo_id> Todo not found").WithLevel(LevelWarning))
		http.Error(w, fmt.Sprintf("Todo with id %d not found", id), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Println(Log("GET todo/<todo_id>"))
		json.NewEncoder(w).Encode(todo)
	case http.MethodPut:
		var updatedTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			log.Println(Log("PUT todo/{todo_id} Invalid item format").WithLevel(LevelError))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updatedTodo.Item == "" {
			log.Println(Log("PUT todo/{todo_id} Item cannot be empty").WithLevel(LevelWarning))
			http.Error(w, "Item cannot be empty", http.StatusBadRequest)
			return
		}

		todo.Item = updatedTodo.Item
		log.Println(Log("PUT todo/{todo_id}"))
		json.NewEncoder(w).Encode(todo)
	case http.MethodDelete:
		for idx := range TodoList {
			if TodoList[idx].ID == id {
				mu.Lock()
				TodoList = append(TodoList[:idx], TodoList[idx+1:]...)
				mu.Unlock()
				log.Println(Log("DELETE todo/<todo_id>"))
				json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
				return
			}
		}
	default:
		log.Println(Log("todo/<todo_id> Method not allowed").WithLevel(LevelWarning))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
