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
	"survey-api/pkg/poll/api/pollid"
	"survey-api/pkg/poll/api/polls"
	pollvote "survey-api/pkg/poll/api/vote"
	"survey-api/pkg/user/api/userid"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/register", register.Handler())
	router.HandleFunc("/login", login.Handler())
	router.HandleFunc("/logout", logout.Handler())
	router.HandleFunc("/token/refresh", refresh.Handler())
	router.HandleFunc("/polls", polls.Handler())
	router.HandleFunc(pollvote.ApiPath, pollvote.Handler())
	router.HandleFunc(pollid.ApiPath, pollid.Handler())
	router.HandleFunc(userid.ApiPath, userid.Handler())

	err := http.ListenAndServe(host+":"+port, router)
	if err != nil {
		panic(err)
	}

	log.Println("Survey server is running on http://" + host + ":" + port)
}
