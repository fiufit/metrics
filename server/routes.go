package server

import "github.com/fiufit/metrics/middleware"
import _ "github.com/fiufit/metrics/docs"
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files"       // swagger embed files

func (s *Server) InitRoutes() {
	baseRouter := s.router.Group("/:version")
	metricsRouter := baseRouter.Group("/metrics")

	baseRouter.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	metricsRouter.POST("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.createMetric.Handle(),
	}))
	metricsRouter.GET("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getMetrics.Handle(),
	}))
}
