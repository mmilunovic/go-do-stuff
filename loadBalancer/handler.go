package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func AllUsers(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "../test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	respondJSON(w, http.StatusOK, users)
}

func NewUser(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "../test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]
	pass := vars["password"]

	user := User{}
	user.Name = name
	user.Email = email
	user.Password = pass

	db.Create(&user)
	respondJSON(w, http.StatusCreated, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	db, err = gorm.Open("sqlite3", "../test.db")
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
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
