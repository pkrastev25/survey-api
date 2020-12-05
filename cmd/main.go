// This package is used only for DEV purposes. It is NOT used for PROD.
//
// The application is deployed as multiple serverless functions.
// The hosting provider looks for function that comply with the
// following signature:
// func(http.ResponseWriter, *http.Request)
//
// For further deployment details, refer to now.json.
// The file contains the routing definitions, which should be an
// exact match with the routing in this file.
package main

import (
	"log"
	"net/http"
	"os"
	"survey-api/pkg/auth/api/login"
	"survey-api/pkg/auth/api/logout"
	"survey-api/pkg/auth/api/refresh"
	"survey-api/pkg/auth/api/register"
	pollapi "survey-api/pkg/poll/api"
	pollvote "survey-api/pkg/poll/api/vote"
)

func main() {
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	http.HandleFunc("/register", register.Handler())
	http.HandleFunc("/login", login.Handler())
	http.HandleFunc("/logout", logout.Handler())
	http.HandleFunc("/token/refresh", refresh.Handler())
	http.HandleFunc("/poll", pollapi.Handler())
	http.HandleFunc("/poll/vote", pollvote.Handler())

	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Survey server is running on http://" + host + ":" + port)
}
