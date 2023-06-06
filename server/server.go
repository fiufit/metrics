package server

import (
	"context"
	"fmt"
	"os"

	"github.com/fiufit/metrics/database"
	"github.com/fiufit/metrics/handlers"
	"github.com/fiufit/metrics/queue_services"
	"github.com/fiufit/metrics/repositories"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server struct {
	router          *gin.Engine
	db              *mongo.Client
	queueConn       *amqp.Connection
	createMetric    handlers.CreateMetric
	getMetrics      handlers.GetMetrics
	metricsConsumer queue_services.MetricsConsumer
}

func (s *Server) Run() {
	go s.metricsConsumer.ConsumeMetrics(context.Background())
	err := s.router.Run(fmt.Sprintf("0.0.0.0:%v", os.Getenv("SERVICE_PORT")))
	if err != nil {
		panic(err)
	}
	if err := s.db.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	db := database.NewMongoDBClient()
	logger, _ := zap.NewDevelopment()

	queueUrl := os.Getenv("CLOUDAMQP_URL")
	if queueUrl == "" {
		queueUrl = "amqp://localhost"
	}
	queueTopic := os.Getenv("METRICS_QUEUE_TOPIC")
	if queueTopic == "" {
		queueTopic = "amqp.fiufitmetrics"
	}
	queueConnection, err := amqp.Dial(queueUrl)
	if err != nil {
		panic(err)
	}

	metricsRepo := repositories.NewMetricsRepository(db, database.DbName, logger)
	metricsPublisher := queue_services.NewMetricsPublisher(queueConnection, queueTopic)
	metricsConsumer := queue_services.NewMetricsConsumer(queueConnection, queueTopic, metricsRepo, logger)

	getMetrics := handlers.NewGetMetrics(metricsRepo)
	createMetric := handlers.NewCreateMetric(metricsPublisher)

	return &Server{
		router:          gin.Default(),
		db:              db,
		queueConn:       queueConnection,
		getMetrics:      getMetrics,
		createMetric:    createMetric,
		metricsConsumer: metricsConsumer,
	}
}
