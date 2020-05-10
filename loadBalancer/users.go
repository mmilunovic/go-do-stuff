package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	pass := vars["password"]

	user := User{}
	user.Name = name
	user.Email = email
	user.Password = pass

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

func getUserByName(name string, w http.ResponseWriter, r *http.Request) *User {
	user := User{}
	if err := db.First(&user, User{Name: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error.")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

var mySigningKey = []byte("mysuperpassword")

func GenerateJWT(w http.ResponseWriter, username string) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	exp := time.Now().Add(time.Minute * 2000)

	claims["authorized"] = true
	claims["user"] = "milossmi"
	claims["exp"] = exp.Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong, %s", err.Error())
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Token",
		Value:   tokenString,
		Expires: exp,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	pass := vars["password"]

	user := getUserByName(name, w, r)

	if user == nil {
		fmt.Fprintf(w, "User not registered")
	}
	if user.Password != pass {
		fmt.Fprintf(w, "Password is incorect.")
	} else {
		GenerateJWT(w, name)
		fmt.Fprintf(w, "Logged In")
	}
}
