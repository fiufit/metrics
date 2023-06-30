# metrics
Metrics microservice for serving backoffice data visualization

[![codecov](https://codecov.io/github/fiufit/metrics/branch/main/graph/badge.svg?token=3QE1J6OCC2)](https://codecov.io/github/fiufit/metrics)
[![test](https://github.com/fiufit/metrics/actions/workflows/test.yml/badge.svg)](https://github.com/fiufit/metrics/actions/workflows/test.yml)
[![Lint Go Code](https://github.com/fiufit/metrics/actions/workflows/lint.yml/badge.svg)](https://github.com/fiufit/metrics/actions/workflows/lint.yml)
[![Fly Deploy](https://github.com/fiufit/metrics/actions/workflows/fly.yml/badge.svg)](https://github.com/fiufit/metrics/actions/workflows/fly.yml)

### Usage

#### With docker:
* Edit .example-env with your own secret credentials and rename it to .env
* `docker build -t fiufit-metrics .`
* `docker run -p PORT:PORT --env-file=.env fiufit-metrics`

#### Natively: 
* `go mod tidy`
* set your environvent variables to imitate the .env-example
* `go run main.go` or `go build` and run the executable


#### Running tests:
* `go test ./...`


### Links
* Fly.io deploy dashboard: https://fly.io/apps/fiufit-metrics
* Swagger docs: https://fiufit-metrics.fly.dev/v1/docs/index.html
* Coverage report: https://app.codecov.io/github/fiufit/metrics
