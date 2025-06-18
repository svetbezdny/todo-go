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

	mux := http.NewServeMux()

	mux.HandleFunc("/todo", TodoHandler)
	mux.HandleFunc("/todo/", TodoByIdHandler)

	handler := LogMiddleware(mux)

	err := http.ListenAndServe(":3000", handler)
	if err != nil {
		panic(err)
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST, DELETE")
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {

	case http.MethodGet:
		json.NewEncoder(w).Encode(TodoList)

	case http.MethodPost:
		var newTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if newTodo.Item == "" {
			http.Error(w, "Item cannot be empty", http.StatusBadRequest)
			return
		}

		todo := Todo{ID: nexId, Item: newTodo.Item}
		nexId++
		mu.Lock()
		TodoList = append(TodoList, todo)
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	case http.MethodDelete:
		TodoList = []Todo{}
		json.NewEncoder(w).Encode(map[string]string{"message": "Todo cleared successfully"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func TodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	todoId := strings.TrimPrefix(r.URL.Path, "/todo/")
	id, err := strconv.Atoi(todoId)
	if err != nil {
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
		http.Error(w, fmt.Sprintf("Todo with id %d not found", id), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(todo)
	case http.MethodPut:
		var updatedTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updatedTodo.Item == "" {
			http.Error(w, "Item cannot be empty", http.StatusBadRequest)
			return
		}

		todo.Item = updatedTodo.Item
		json.NewEncoder(w).Encode(todo)
	case http.MethodDelete:
		for idx := range TodoList {
			if TodoList[idx].ID == id {
				mu.Lock()
				TodoList = append(TodoList[:idx], TodoList[idx+1:]...)
				mu.Unlock()
				json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
				return
			}
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := &LogResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		handler.ServeHTTP(mw, r)
		log.Printf("%s [%d]", r.URL.RequestURI(), mw.StatusCode)
	})
}
