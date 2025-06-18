package main

import (
	"net/http"
)

type Todo struct {
	ID   int    `json:"id"`
	Item string `json:"item"`
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
