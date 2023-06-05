package testing

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestSuite struct {
	DB       *mongo.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func NewTestSuite() TestSuite {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			"MONGO_INITDB_DATABASE=testdb",
			"MONGO_INITDB_ROOT_USERNAME=testmongouser",
			"MONGO_INITDB_ROOT_PASSWORD=testmongopassword",
		},
	},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	var mongoClient *mongo.Client
	if err = pool.Retry(func() error {
		mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(
			fmt.Sprintf("mongodb://testmongouser:testmongopassword@localhost:%s", resource.GetPort("27017/tcp")),
		))
		if err != nil {
			return err
		}

		return mongoClient.Database("testdb").RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	suite := TestSuite{
		DB:       mongoClient,
		pool:     pool,
		resource: resource,
	}

	return suite
}

func (ts TestSuite) TearDown() {
	if err := ts.pool.Purge(ts.resource); err != nil {
		log.Fatalf("Could not TearDown testing DB: %s", err)
	}

	if err := ts.DB.Disconnect(context.Background()); err != nil {
		log.Fatalf("Could not disconnect from testing DB: %s", err)
	}
}

func (ts TestSuite) Truncate(dbName string, collection string) {
	_ = ts.DB.Database(dbName).Collection(collection).Drop(context.Background())
}
