package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/lb/tasks/{userID}", AllTasks).Methods("GET")
	myRouter.HandleFunc("/lb/task/{userID}/{title}", NewTask).Methods("POST")
	myRouter.HandleFunc("/lb/task/{userID}/{taskID}", DeleteTask).Methods("DELETE")
	myRouter.HandleFunc("/lb/task/update/{userID}/{taskID}", UpdateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8086", myRouter))
}

func main() {

	InitialMigration()

	handleRequests()
}
