package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	ErrNotFound = errors.New("not found")
)

func (Todo) TableName() string {
	return "todos"
}

func InitDatabase() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	DB.AutoMigrate(&Todo{})
	return nil
}

func InsertTodo(db *gorm.DB, todo Todo) (Todo, error) {
	err := db.Create(&todo).Error
	return todo, err
}

func GetTodoById(db *gorm.DB, id int) (Todo, error) {
	var todo Todo
	err := db.Where("id = ?", id).Take(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Todo{}, fmt.Errorf("Todo with id %d not found", id)
		}
	}
	return todo, err
}

func UpdateTodoById(db *gorm.DB, id int, item string) (Todo, error) {
	var todo Todo
	err := db.Model(&Todo{}).Where("id = ?", id).Update("item", item).First(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Todo{}, fmt.Errorf("Todo with id %d not found", id)
		}
	}
	return todo, err
}

func DeleteTodoById(db *gorm.DB, id int) (bool, error) {
	result := db.Where("id = ?", id).Delete(&Todo{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func GetAllTodo(db *gorm.DB) ([]Todo, error) {
	todos := make([]Todo, 0)
	err := db.Find(&todos).Error
	return todos, err
}

func DeleteAllTodo(db *gorm.DB) error {
	result := db.Where("1 = 1").Delete(&Todo{})
	return result.Error
}
