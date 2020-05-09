package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tenderly")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")

	// HandleFunc User CRUD operations
	myRouter.HandleFunc("/users", AllUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}/{email}", NewUser).Methods("POST")
	myRouter.HandleFunc("/user/{name}", DeleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", UpdateUser).Methods("PUT")

	// HandleFunc Task CRUD operations
	myRouter.HandleFunc("/tasks/{userID}", AllTasks).Methods("GET")
	myRouter.HandleFunc("/task/{userID}/{title}", NewTask).Methods("POST")
	myRouter.HandleFunc("/task/{userID}/{taskID}", DeleteTask).Methods("DELETE")
	myRouter.HandleFunc("/task/update/{userID}/{taskID}", UpdateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	InitialMigration()

	handleRequests()
}
