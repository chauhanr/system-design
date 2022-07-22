package app

import "github.com/gorilla/mux"

func (s *server) ConfigureRoutes() {
	var api = s.Router.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = s.NotFoundHanlder()
	s.apiV1Routes(api)
}

func (s *server) apiV1Routes(subRouter *mux.Router) {
	var api = subRouter.PathPrefix("/v1").Subrouter()
	api.NotFoundHandler = s.NotFoundHanlder()
	api.HandleFunc("/{domain}/{userId}", s.ForwardingHandler())
}
