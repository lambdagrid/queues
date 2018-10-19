package server

import (
	"github.com/lambdagrid/queues/auth"
	"github.com/lambdagrid/queues/middleware"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	authProvider auth.AuthProvider
	router       *httprouter.Router
}

func New(authProvider auth.AuthProvider) Server {
	s := Server{
		authProvider: authProvider,
		router:       httprouter.New(),
	}
	s.routes()

	return s
}

func (s Server) GetRouter() *httprouter.Router {
	return s.router
}

func (s Server) routes() {
	s.router.GET("/status", s.getStatus())
	s.router.POST("/jobs", middleware.HeaderAuth(s.authProvider, s.handleJobWrite()))
}
