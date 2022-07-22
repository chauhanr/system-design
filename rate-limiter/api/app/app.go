package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	Router  *mux.Router
	hClient *http.Client
}

func NewServer(router *mux.Router, hClient *http.Client) *server {
	s := server{Router: router, hClient: hClient}
	return &s
}

func Startup(hClient *http.Client) {
	r := mux.NewRouter()
	s := NewServer(r, hClient)

	s.ConfigureRoutes()

	// TODO. replace the 8080 by using viper and env files.
	err := http.ListenAndServe(":8080", s.Router)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}
