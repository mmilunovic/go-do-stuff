package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var workerList WorkerList
var myRouter = mux.NewRouter().StrictSlash(true)

func handleRequests() {

	workerRoutes := []string{
		"localhost:8083",
		"localhost:8086",
		"localhost:8088",
	}

	workerList.init(workerRoutes)

	/*
		This should use Regex or something but I can't figure it out...
	*/
	myRouter.Handle("/tasks/{userID}", IsAuthorized(func(w http.ResponseWriter, r *http.Request) {
		workerList.loadBalance(w, r)
	}))
	myRouter.Handle("/task/{userID}/{title}", IsAuthorized(func(w http.ResponseWriter, r *http.Request) {
		workerList.loadBalance(w, r)
	}))
	myRouter.Handle("/task/{userID}/{taskID}", IsAuthorized(func(w http.ResponseWriter, r *http.Request) {
		workerList.loadBalance(w, r)
	}))
	myRouter.Handle("/task/update/{userID}/{taskID}", IsAuthorized(func(w http.ResponseWriter, r *http.Request) {
		workerList.loadBalance(w, r)
	}))

	myRouter.HandleFunc("/login/{name}/{password}", LoginHandler).Methods("POST")
	myRouter.HandleFunc("/register/{name}/{email}/{password}", NewUser).Methods("POST")

	myRouter.HandleFunc("/users", AllUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", DeleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{userID}", UpdateUser).Methods("PUT")

	http.Handle("/", myRouter)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {

	InitialMigration()

	handleRequests()
}
