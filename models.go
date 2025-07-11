package main

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Item string `json:"item" validate:"required,min=1"`
}

type TodoResponse struct {
	Item      string    `json:"item"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type TodoCreateResponse struct {
	Item      string    `json:"item"`
	CreatedAt time.Time `json:"createdAt"`
}

type TodoUpdateResponse struct {
	Item      string    `json:"item"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TodoItem struct {
	Item string `json:"item"`
}

func (Todo) TableName() string {
	return "todos"
}

type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *LogResponseWriter) WriteHeader(StatusCode int) {
	w.ResponseWriter.WriteHeader(StatusCode)
	w.StatusCode = StatusCode
}
