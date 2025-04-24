package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"grpc.bank-api/internal/config"
	"grpc.bank-api/internal/db"
)

type Server struct {
	logger *zap.Logger
	config *config.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config *config.Config, store db.Store, log *zap.Logger) *Server {
	server := &Server{
		config: config,
		store:  store,
		router: gin.Default(),
		logger: log,
	}
	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	router := gin.Default()
	api := router.Group("/api/v1")

	api.GET("/health", s.healthCheck)

	api.GET("/todos", s.getAllTodos)
	api.POST("/todos", s.createTodo)
	api.GET("/todos/:id", s.getTodo)
	api.PUT("/todos/:id", s.completeTodo)
	api.PUT("/todos/:id", s.updateTodo)
	api.DELETE("/todos/:id", s.deleteTodo)

	s.router = router
}

func (s *Server) Start() error {
	s.logger.Info("starting server", zap.String("address", s.config.HttpServerAddress))

	return s.router.Run(s.config.HttpServerAddress)
}

func errResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message,
	})
}
