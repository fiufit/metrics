package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/fiufit/metrics/contracts/metrics"
	"github.com/fiufit/metrics/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap/zaptest"
)

func TestMetricsRepository_Create_Ok(t *testing.T) {
	defer testSuite.Truncate(metricsCollection)
	ctx := context.Background()

	db := testSuite.DB
	repo := NewMetricsRepository(db, testDbName, zaptest.NewLogger(t))
	testMetric := models.Metric{
		MetricType: "login",
		SubType:    "mail",
		DateTime:   time.Now(),
	}

	_ = repo.Create(ctx, testMetric)

	var createdMetric models.Metric
	err := repo.db.FindOne(ctx, bson.D{}).Decode(&createdMetric)

	assert.NoError(t, err)
	assert.Equal(t, testMetric.MetricType, createdMetric.MetricType)
}

func TestMetricsRepository_Get_Ok(t *testing.T) {
	defer testSuite.Truncate(metricsCollection)
	ctx := context.Background()
	db := testSuite.DB
	repo := NewMetricsRepository(db, testDbName, zaptest.NewLogger(t))

	testMetrics := []interface{}{
		models.Metric{MetricType: "login", SubType: "mail", DateTime: time.Now()},
		models.Metric{MetricType: "login", SubType: "federated_entity", DateTime: time.Now()},
		models.Metric{MetricType: "register", SubType: "mail", DateTime: time.Now().AddDate(-1, 0, 0)},
		models.Metric{MetricType: "register", SubType: "mail", DateTime: time.Now().AddDate(-1, 0, 0)},
		models.Metric{MetricType: "register", SubType: "federated_entity", DateTime: time.Now().AddDate(1, 0, 0)},
		models.Metric{MetricType: "blocked", SubType: "", DateTime: time.Now()},
		models.Metric{MetricType: "password_recover", SubType: "", DateTime: time.Now()},
		models.Metric{MetricType: "password_recover", SubType: "", DateTime: time.Now()},
		models.Metric{MetricType: "location", SubType: "", DateTime: time.Now()},
		models.Metric{MetricType: "new_training", SubType: "", DateTime: time.Now()},
	}

	_, _ = repo.db.InsertMany(ctx, testMetrics)

	type testCase struct {
		description      string
		expectedResCount int
		req              metrics.GetMetricsRequest
	}

	for _, tCase := range []testCase{
		{
			description:      "login type any subtype",
			expectedResCount: 2,
			req:              metrics.GetMetricsRequest{MetricType: "login", SubType: "", From: time.Time{}, To: time.Time{}},
		},
		{
			description:      "login type mail subtype",
			expectedResCount: 1,
			req:              metrics.GetMetricsRequest{MetricType: "login", SubType: "mail", From: time.Time{}, To: time.Time{}},
		},
		{
			description:      "register type any date",
			expectedResCount: 3,
			req:              metrics.GetMetricsRequest{MetricType: "register", SubType: "", From: time.Time{}, To: time.Time{}},
		},
		{
			description:      "register type year-old dates",
			expectedResCount: 2,
			req:              metrics.GetMetricsRequest{MetricType: "register", SubType: "", From: time.Now().AddDate(-1, 0, -1), To: time.Now().AddDate(0, 0, -1)},
		},
		{
			description:      "blocked type",
			expectedResCount: 1,
			req:              metrics.GetMetricsRequest{MetricType: "blocked", SubType: "", From: time.Time{}, To: time.Time{}},
		},
		{
			description:      "password_recover type",
			expectedResCount: 2,
			req:              metrics.GetMetricsRequest{MetricType: "password_recover", SubType: "", From: time.Time{}, To: time.Time{}},
		},
		{
			description:      "location type",
			expectedResCount: 1,
			req:              metrics.GetMetricsRequest{MetricType: "location", SubType: "", From: time.Time{}, To: time.Time{}},
		},
		{
			description:      "new_training type",
			expectedResCount: 1,
			req:              metrics.GetMetricsRequest{MetricType: "new_training", SubType: "", From: time.Time{}, To: time.Time{}},
		},
	} {
		t.Run(tCase.description, func(t *testing.T) {
			res, err := repo.Get(ctx, tCase.req)
			assert.NoError(t, err)
			assert.Equal(t, tCase.expectedResCount, len(res))
		})
	}
}
