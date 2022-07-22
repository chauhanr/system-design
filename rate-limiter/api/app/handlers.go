package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) NotFoundHanlder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *server) ForwardingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathVars := mux.Vars(r)
		userId := pathVars["userId"]
		domain := pathVars["domain"]
		log.Printf("Domain: %s, UserId: %s\n", domain, userId)

	}
}
