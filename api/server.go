package api

import (
	db "GoBankProject/db/sqlc"
	"github.com/gin-gonic/gin"
	"time"
)

// Server will serve all HTTP requests for our banking services
type Server struct {
	store *db.Store
	router *gin.Engine
}

// NewServer will create a new server instance and set up all api routes for our services on that server
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	//	add routes to router
	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)

	router.GET("/accounts", server.listAccount)
	server.router = router
	return server
}
func (s *Server) Run(address string) error {
	return s.router.Run(address)
}


func errorResponse(status int, err error) gin.H {
	// TODO: convert the error to a better format
	return gin.H{
		"status": status,
		"error": err.Error(),
		"time": time.Now(),
	}
}
