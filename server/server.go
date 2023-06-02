package server

import (
	"context"
	"fmt"
	"os"

	"github.com/fiufit/metrics/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router *gin.Engine
	db     *mongo.Client
}

func (s *Server) Run() {
	err := s.router.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("SERVICE_PORT")))
	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	db := database.NewMongoDBClient()
	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	return &Server{
		router: gin.Default(),
		db:     db,
	}
}
