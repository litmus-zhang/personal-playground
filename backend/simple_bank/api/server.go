package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/litmus-zhang/simple_bank/bank"
)

type Server struct {
	store  bank.Store
	router *gin.Engine
}

func NewServer(store bank.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	api := router.Group("/api/v1")

	api.POST("/accounts", server.createAccount)
	api.GET("/accounts/:id", server.getAccount)
	api.GET("/accounts", server.listAccounts)
	api.POST("/transfers", server.createTransfer)
	api.POST("/users", server.createUser)
	api.GET("/users/:username", server.getUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
