package main

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type DBRepository struct {
	db *gorm.DB
}

type Repository interface {
	Insert(todo Todo) (Todo, error)
	GetById(id int) (Todo, error)
	UpdateById(id int, item string) (Todo, error)
	DeleteById(id int) error
	GetAll() ([]Todo, error)
	DeleteAll() error
}

func NewDBRepository() (*DBRepository, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DBRepository{db: db}, nil
}

func (db *DBRepository) InitDatabase() error {
	if err := db.db.AutoMigrate(&Todo{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}

func (db *DBRepository) Insert(todo Todo) (Todo, error) {
	if err := db.db.Create(&todo).Error; err != nil {
		return Todo{}, fmt.Errorf("failed to insert todo: %w", err)
	}
	return todo, nil
}

func (db *DBRepository) GetById(id int) (Todo, error) {
	var todo Todo
	err := db.db.Where("id = ?", id).Take(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Todo{}, fmt.Errorf("todo with id %d not found", id)
		}
		return Todo{}, fmt.Errorf("failed to get todo by id: %w", err)
	}
	return todo, nil
}

func (db *DBRepository) UpdateById(id int, item string) (Todo, error) {
	var todo Todo
	err := db.db.Model(&Todo{}).Where("id = ?", id).Update("item", item).First(&todo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Todo{}, ErrNotFound
		}
		return Todo{}, fmt.Errorf("failed to update todo by id: %w", err)
	}
	return todo, nil
}

func (db *DBRepository) DeleteById(id int) error {
	result := db.db.Where("id = ?", id).Delete(&Todo{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete todo by id: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *DBRepository) GetAll() ([]Todo, error) {
	var todos []Todo
	if err := db.db.Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to get all todos: %w", err)
	}
	return todos, nil
}

func (db *DBRepository) DeleteAll() error {
	if err := db.db.Where("1 = 1").Delete(&Todo{}).Error; err != nil {
		return fmt.Errorf("failed to delete all todos: %w", err)
	}
	return nil
}
