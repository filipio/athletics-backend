package controllers

import (
	"fmt"
	"net/http"
	"strconv"
)

func HandleSomething() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, World!")
		},
	)
}

func HandleSomePut() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Error(w, "Invalid id", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Put request with id %d", id)
		},
	)
}

func HandleDeepPut() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Error(w, "Invalid id", http.StatusBadRequest)
				return
			}

			otherId, err := strconv.Atoi(r.PathValue("otherId"))
			if err != nil {
				http.Error(w, "Invalid otherId", http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, "Put request with id %d and otherId %d", id, otherId)
		},
	)
}
