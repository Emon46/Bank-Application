package api

import (
	db "github.com/emon46/bank-application/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer it will create a new gin api and setup routing for all the api call
func NewServer(store db.Store) *Server {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validateCurrency)
		if err != nil {
			log.Fatal(err)
		}
	}
	server := &Server{
		store: store,
	}
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/accounts/:id", server.updateAccount)
	router.POST("/transfer", server.createTransfer)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
