package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

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

func GenerateJWT(w http.ResponseWriter, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	exp := time.Now().Add(time.Minute * 2000)

	claims["authorized"] = true
	claims["user"] = username
	claims["exp"] = exp.Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong, %s", err.Error())
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Token",
		Value:   tokenString,
		Expires: exp,
	})
	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("sqlite3", "../test.db")
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
		validToken, _ := GenerateJWT(w, name)
		if validToken != "" {
			fmt.Fprintf(w, validToken)
			respondJSON(w, http.StatusOK, nil)
		} else {
			fmt.Fprintf(w, "Token is null for some reason")
		}
	}
}
