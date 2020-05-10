package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/tasks/{userID}", AllTasks).Methods("GET")
	myRouter.HandleFunc("/task/{userID}/{title}", NewTask).Methods("POST")
	myRouter.HandleFunc("/task/{userID}/{taskID}", DeleteTask).Methods("DELETE")
	myRouter.HandleFunc("/task/update/{userID}/{taskID}", UpdateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8086", myRouter))
}

func main() {

	InitialMigration()

	handleRequests()
}
