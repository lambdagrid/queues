package server

import (
	"github.com/jmoiron/sqlx"
	"github.com/lambdagrid/queues/auth"
	"github.com/lambdagrid/queues/middleware"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	authProvider auth.AuthProvider
	router       *httprouter.Router
	DB           *sqlx.DB // TODO: Abstract much of the functionality away
}

func New(authProvider auth.AuthProvider, db *sqlx.DB) Server {
	s := Server{
		authProvider: authProvider,
		router:       httprouter.New(),
		DB:           db,
	}
	s.routes()

	return s
}

func (s Server) GetRouter() *httprouter.Router {
	return s.router
}

func (s Server) routes() {
	s.router.GET("/v1/status", s.getStatus())
	s.router.POST("/v1/jobs", middleware.HeaderAuth(s.authProvider, s.handleJobWrite()))
	s.router.POST("/v1/signup", s.signup())
	s.router.GET("/v1/queue/list", middleware.HeaderAuth(s.authProvider, s.listQueues()))
	s.router.POST("/v1/queue/create", middleware.HeaderAuth(s.authProvider, s.createQueue()))
	// queueID will be an account scoped unique name
	s.router.GET("/v1/queue/info/:queuename", middleware.HeaderAuth(s.authProvider, s.queueInfo()))
	// put queue name in payload
	s.router.POST("/v1/queue/message/send", middleware.HeaderAuth(s.authProvider, s.sendMessage()))
	s.router.GET("/v1/queue/message/receive/:queuename", middleware.HeaderAuth(s.authProvider, s.receiveMessage()))
	// queue name in payload
	s.router.POST("/v1/queue/message/complete", middleware.HeaderAuth(s.authProvider, s.completeMessage()))
	s.router.POST("/v1/queue/message/fail", middleware.HeaderAuth(s.authProvider, s.completeMessage()))
}
