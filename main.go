package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	if err := InitDatabase(); err != nil {
		log.Fatal("Failed to connect to database")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/todo", TodoHandler)
	mux.HandleFunc("/todo/", TodoByIdHandler)
	handler := LogMiddleware(mux)

	log.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST, DELETE")
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		todos, err := GetAllTodo(DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(todos)

	case http.MethodPost:
		var newTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := InsertTodo(DB, Todo{Item: newTodo.Item})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)

	case http.MethodDelete:
		if err := DeleteAllTodo(DB); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Todo cleared successfully"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func TodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	todoId := strings.TrimPrefix(r.URL.Path, "/todo/")

	id, err := strconv.Atoi(todoId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid todo id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, err := GetTodoById(DB, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(todo)
	case http.MethodPut:
		var updatedTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := UpdateTodoById(DB, id, updatedTodo.Item)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(todo)
	case http.MethodDelete:
		ok, err := DeleteTodoById(DB, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Todo with id %d not found", id)})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		mw := &LogResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		handler.ServeHTTP(mw, r)
		log.Printf(`{
			"host":"%s",
			"method":"%s",
			"path":"%s",
			"status":%d,
			"duration":"%s"
			}`,
			r.RemoteAddr, r.Method, r.URL.Path, mw.StatusCode, time.Since(start).String())
	})
}
