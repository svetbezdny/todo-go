package main

type Todo struct {
	ID   int    `json:"id"`
	Item string `json:"item"`
}
type TodoItem struct {
	Item string `json:"item"`
}
