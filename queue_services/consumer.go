package queue_services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fiufit/metrics/contracts/metrics"
	"github.com/fiufit/metrics/models"
	"github.com/fiufit/metrics/repositories"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type MetricsConsumer struct {
	queueConn *amqp.Connection
	topic     string
	metrics   repositories.Metrics
	logger    *zap.Logger
}

func NewMetricsConsumer(queueConn *amqp.Connection, topic string, metrics repositories.Metrics, logger *zap.Logger) MetricsConsumer {
	return MetricsConsumer{queueConn: queueConn, topic: topic, metrics: metrics, logger: logger}
}

func (mc MetricsConsumer) ConsumeMetrics(ctx context.Context) {
	channel, err := mc.queueConn.Channel()
	if err != nil {
		mc.logger.Error("Unable to open amqp channel", zap.Error(err))
		return
	}

	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			mc.logger.Error("Unable to close amqp channel", zap.Error(err))
		}
	}(channel)

	queue, err := channel.QueueDeclare("metrics", false, true, false, true, nil)
	if err != nil {
		mc.logger.Error("Unable to declare amqp queue", zap.Error(err))
		return
	}

	err = channel.QueueBind(queue.Name, "metric", mc.topic, false, nil)
	if err != nil {
		mc.logger.Error("Unable to bind amqp queue", zap.Error(err))
		return
	}

	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		mc.logger.Error("Unable to consume amqp channel", zap.Error(err))
		return
	}

	func(messages <-chan amqp.Delivery) {

		for {
			select {
			case msg := <-messages:
				var req metrics.CreateMetricRequest
				err := json.Unmarshal(msg.Body, &req)

				if err != nil {
					mc.logger.Error("Unable to decode consumed metric", zap.Error(err))
					continue
				}

				metric := models.Metric{
					MetricType: req.MetricType,
					SubType:    req.SubType,
					DateTime:   time.Now(),
				}

				_ = mc.metrics.Create(ctx, metric)

				if channel.IsClosed() {
					mc.logger.Info("Metrics channel has been closed, stopping consumption")
					return
				}

			case <-ctx.Done():
				mc.logger.Info("Metrics consumption has been context-cancelled")
				return
			}
		}
	}(messages)
}
