package main

import (
	"fmt"
	"net/http/httputil"

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

type Worker struct {
	Route        string
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type WorkerList struct {
	Workers []Worker
	Latest  int
}

func InitialMigration() {
	db, err = gorm.Open("sqlite3", "../test.db")

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database.")
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Task{})
	db.Model(&Task{}).AddForeignKey("userID", "users(id)", "CASCADE", "CASCADE")
}
