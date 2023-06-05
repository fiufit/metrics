package repositories

import (
	"context"

	mContracts "github.com/fiufit/metrics/contracts/metrics"
	"github.com/fiufit/metrics/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const metricsCollection = "metrics"

type Metrics interface {
	Get(ctx context.Context, req mContracts.GetMetricsRequest) (mContracts.GetMetricsResponse, error)
	Create(ctx context.Context, req mContracts.GetMetricsResponse) error
}

type MetricsRepository struct {
	logger *zap.Logger
	db     *mongo.Collection
}

func NewMetricsRepository(db *mongo.Client, dbName string, logger *zap.Logger) MetricsRepository {
	return MetricsRepository{db: db.Database(dbName).Collection(metricsCollection), logger: logger}
}

func (repo MetricsRepository) Get(ctx context.Context, req mContracts.GetMetricsRequest) (mContracts.GetMetricsResponse, error) {

	filter := bson.M{}
	filter["type"] = req.MetricType

	if req.SubType != "" {
		filter["subtype"] = req.SubType
	}
	if !req.From.IsZero() && !req.To.IsZero() {
		filter["timestamp"] = bson.M{
			"$gte": req.From,
			"$lt":  req.To,
		}
	}

	sort := bson.D{{Key: "timestamp", Value: 1}}
	opts := options.Find().SetSort(sort)

	cursor, err := repo.db.Find(ctx, filter, opts)
	if err != nil {
		repo.logger.Error("Unable to query metrics from collection", zap.Error(err))
		return mContracts.GetMetricsResponse{}, err
	}

	var metrics []models.Metric
	if err = cursor.All(ctx, &metrics); err != nil {
		repo.logger.Error("Unable to unmarshall metrics from collection", zap.Error(err))
		return mContracts.GetMetricsResponse{}, err
	}

	return metrics, nil
}

func (repo MetricsRepository) Create(ctx context.Context, metric models.Metric) error {
	_, err := repo.db.InsertOne(ctx, metric)

	if err != nil {
		repo.logger.Error("Failed to insert metric in metrics collection", zap.Error(err))
	}
	return err
}
