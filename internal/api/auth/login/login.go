package login

import (
	"net/http"
)

var handler func(http.ResponseWriter, *http.Request)

func Handler() func(http.ResponseWriter, *http.Request) {
	return handler
}

func Init() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte("Login called!"))
		w.WriteHeader(http.StatusOK)
	}
}

func init() {
	handler = Init()
}
