package server

import "github.com/fiufit/metrics/middleware"

func (s *Server) InitRoutes() {
	baseRouter := s.router.Group("/:version")
	metricsRouter := baseRouter.Group("/metrics")

	metricsRouter.POST("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.createMetric.Handle(),
	}))
	metricsRouter.GET("", middleware.HandleByVersion(middleware.VersionHandlers{
		"v1": s.getMetrics.Handle(),
	}))
}
