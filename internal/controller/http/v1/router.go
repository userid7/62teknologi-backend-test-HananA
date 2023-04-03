// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"backend-test/internal/business/usecase"
	"backend-test/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, b usecase.Business) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// // Swagger
	// swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	// handler.GET("/swagger/*any", swaggerHandler)

	// // K8s probe
	// handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// // Routers
	// h := handler.Group("/bussiness")
	// {
	// 	newTranslationRoutes(h, t, l)
	// }
	h := handler.Group("/")
	newBusinessRoutes(h, b, l)
}
