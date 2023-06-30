package main

import "github.com/fiufit/metrics/server"

//	@title			Fiufit Metrics API
//	@version		dev
//	@description	Fiufit's Metrics service documentation. This API serves Fiufit usage metrics for backoffice visualization.

//	@host		fiufit-metrics.fly.dev
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	srv := server.NewServer()
	srv.InitRoutes()
	srv.Run()
}
