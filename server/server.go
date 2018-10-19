package server

import (
	"net/http"

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

func (s Server) stubHandler() httprouter.Handle {
	return func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {}
}

func (s Server) routes() {
	s.router.GET("/v1/status", s.getStatus())
	s.router.POST("/v1/jobs", middleware.HeaderAuth(s.authProvider, s.handleJobWrite()))
	s.router.GET("/v1/queue/list", s.stubHandler())
	s.router.POST("/v1/queue/create", s.stubHandler())
	// queueID will be an account scoped unique name
	s.router.GET("/v1/queue/info/:queuename", s.stubHandler())
	// put queue name in payload
	s.router.POST("/v1/queue/message/send", s.stubHandler())
	s.router.GET("/v1/queue/message/receive/:queuename", s.stubHandler())
	// queue name in payload
	s.router.POST("/v1/queue/message/complete", s.stubHandler())
	s.router.POST("/v1/queue/message/fail", s.stubHandler())
}
