package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

var mySigningKey = []byte("mysuperpassword")

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "milossmi"
	claims["exp"] = time.Now().Add(time.Minute * 2000).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong, %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func homePage(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// client := &http.Client()
	// req, _ := http.NewRequest("GET", "http://localhost:8081/", nil)
	// req.Header.Set("Token", validToken)

	// res,err := client.Do(req)

	// if err != nil {
	// 	fmt.Fprintf(w, "Error: %s", err.Error())
	// }

	fmt.Fprintf(w, validToken)
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
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

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")

	// Handle User CRUD operations
	myRouter.Handle("/users", isAuthorized(AllUsers)).Methods("GET")
	myRouter.Handle("/user/{name}/{email}", isAuthorized(NewUser)).Methods("POST")
	myRouter.Handle("/user/{name}", isAuthorized(DeleteUser)).Methods("DELETE")
	myRouter.Handle("/user/{name}/{email}", isAuthorized(UpdateUser)).Methods("PUT")

	// Handle Task CRUD operations
	myRouter.Handle("/tasks/{userID}", isAuthorized(AllTasks)).Methods("GET")
	myRouter.Handle("/task/{userID}/{title}", isAuthorized(NewTask)).Methods("POST")
	myRouter.Handle("/task/{userID}/{taskID}", isAuthorized(DeleteTask)).Methods("DELETE")
	myRouter.Handle("/task/update/{userID}/{taskID}", isAuthorized(UpdateTask)).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	InitialMigration()

	handleRequests()
}
