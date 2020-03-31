// This package is used only for DEV purposes. It is NOT used for PROD.
//
// The application is deployed as multiple serverless functions.
// The hosting provider looks for function that comply with the
// following signature:
// (w http.ResponseWriter, r *http.Request)
//
// For further deployed details, refer to now.json.
// The file contains the routing definitions, which should be an
// exact match with the routing in this file.
package main

import (
	"net/http"
	"os"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Survey server is running on http://" + host + ":" + port))
	})
	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		panic(err)
	}
}
