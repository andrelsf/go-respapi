package api

import (
	db "github.com/andrelsf/go-restapi/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.postAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAllAccounts)

	server.router = router
	return server
}

// Start runs the HTTP Server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(code int, err error) gin.H {
	return gin.H{
		"code":    code,
		"message": err.Error(),
	}
}
