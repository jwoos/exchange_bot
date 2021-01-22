package main

func (s *Server) routes() {
	s.router.HandleFunc("/status", s.handleStatus()).Methods("GET")
	s.router.HandleFunc("/events", s.handleEvents()).Methods("POST")
}
