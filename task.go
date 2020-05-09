package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Task struct {
	gorm.Model
	Title  string
	Status bool
	UserID uint
}

func AllTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All Tasks endpoint Hit")
}

func NewTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "New Task endpoint hit.")
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete Task endpoint hit.")
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update Task endpoint hit.")
}
