package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Tasks    []Task `gorm:"ForeignKey:UserID" json:"tasks"`
}

type Task struct {
	gorm.Model
	Title  string
	Status bool
	UserID uint `json:"userID"`
}

func (t *Task) Complete() {
	t.Status = true
}

func (t *Task) Undo() {
	t.Status = false
}

func InitialMigration() {
	db, err = gorm.Open("sqlite3", "test.db")

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database.")
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Task{})
	db.Model(&Task{}).AddForeignKey("userID", "users(id)", "CASCADE", "CASCADE")
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func AllUsers(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	respondJSON(w, http.StatusOK, users)
}

func NewUser(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	user := User{}
	user.Name = name
	user.Email = email

	db.Create(&user)
	respondJSON(w, http.StatusCreated, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	db.Where("name = ?", name).Find(&user)

	db.Delete(&user)

	respondJSON(w, http.StatusNoContent, nil)
	fmt.Fprintf(w, "User successfully deleted.")

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Email = email

	db.Save(&user)

	fmt.Fprintf(w, "User successfully updated.")
}

func getUser(id int, w http.ResponseWriter, r *http.Request) *User {
	user := User{}

	if err := db.First(&user, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}

func getTask(id int, w http.ResponseWriter, r *http.Request) *Task {
	task := Task{}

	if err := db.First(&task, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &task
}

func AllTasks(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
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

	db, err = gorm.Open("sqlite3", "test.db")
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
	db, err = gorm.Open("sqlite3", "test.db")
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
	db, err = gorm.Open("sqlite3", "test.db")
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
