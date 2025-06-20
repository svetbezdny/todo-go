package main

import (
	"net/http"
)

type Todo struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Item string `json:"item" validate:"required,min=1"`
}
type TodoItem struct {
	Item string `json:"item"`
}

type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *LogResponseWriter) WriteHeader(StatusCode int) {
	w.ResponseWriter.WriteHeader(StatusCode)
	w.StatusCode = StatusCode
}
