package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func AllTasks(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "../test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["userID"])

	user := getUser(userID, w, r)

	if user == nil {
		return
	}

	tasks := []Task{}

	if err := db.Model(&user).Related(&tasks).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tasks)

}

func NewTask(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "../test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)

	userID, _ := strconv.Atoi(vars["userID"])
	user := getUser(userID, w, r)
	if user == nil {
		return
	}

	task := Task{UserID: user.ID}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "../test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)

	userID, _ := strconv.Atoi(vars["userID"])
	user := getUser(userID, w, r)
	if user == nil {
		return
	}
	id, _ := strconv.Atoi(vars["taskID"])

	task := getTask(id, w, r)
	if task == nil {
		return
	}

	if err := db.Delete(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "../test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)

	userID, _ := strconv.Atoi(vars["userID"])
	user := getUser(userID, w, r)
	if user == nil {
		return
	}
	id, _ := strconv.Atoi(vars["taskID"])

	task := getTask(id, w, r)
	if task == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&task).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
