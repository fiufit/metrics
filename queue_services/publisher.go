package queue_services

import (
	"context"
	"encoding/json"

	"github.com/fiufit/metrics/contracts/metrics"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MetricsPublisher struct {
	queueConn *amqp.Connection
	topic     string
}

func NewMetricsPublisher(queueConn *amqp.Connection, topic string) MetricsPublisher {
	return MetricsPublisher{queueConn: queueConn, topic: topic}
}

func (mp *MetricsPublisher) PublishMetric(ctx context.Context, req metrics.CreateMetricRequest) error {
	channel, err := mp.queueConn.Channel()
	if err != nil {
		return err
	}

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: 1,
		ContentType:  "text/plain",
		Body:         reqBytes,
	}

	err = channel.PublishWithContext(ctx, mp.topic, "metric", false, false, msg)
	return err
}
