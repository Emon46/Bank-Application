package api

import (
	db "github.com/emon46/bank-application/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer it will create a new gin api and setup routing for all the api call
func NewServer(store *db.Store) *Server {
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
	}
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/accounts/:id", server.updateAccount)

	return &Server{
		store:  store,
		router: router,
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
