package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var repo *DBRepository

func main() {

	var err error
	repo, err = NewDBRepository()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := repo.InitDatabase(); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/todo", TodoHandler)
	mux.HandleFunc("/todo/", TodoByIdHandler)
	handler := LogMiddleware(mux)

	log.Println("starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST, DELETE")
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		todos, err := repo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(todos)

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

		todo, err := repo.Insert(Todo{Item: newTodo.Item})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TodoCreateResponse{Item: todo.Item, CreatedAt: todo.CreatedAt})

	case http.MethodDelete:
		if err := repo.DeleteAll(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "todo cleared successfully"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

func TodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")

	todoId := strings.TrimPrefix(r.URL.Path, "/todo/")

	id, err := strconv.Atoi(todoId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid todo id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, err := repo.GetById(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(TodoResponse{Item: todo.Item, CreatedAt: todo.CreatedAt, UpdatedAt: todo.UpdatedAt})
	case http.MethodPut:
		var updatedTodo TodoItem
		if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo, err := repo.UpdateById(id, updatedTodo.Item)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(TodoUpdateResponse{Item: todo.Item, UpdatedAt: todo.UpdatedAt})
	case http.MethodDelete:
		err := repo.DeleteById(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("todo %d deleted successfully", id)})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
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
