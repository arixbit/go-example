package router

import (
	"github.com/gin-gonic/gin"

	"go-skeleton/internal/handler"
)

// Dependencies collects handlers and middleware needed during route registration.
type Dependencies struct {
	Example *handler.ExampleHandler
}

// RegisterRoutes registers API routes under the given router group.
func RegisterRoutes(r *gin.RouterGroup, deps Dependencies) error {
	registerExampleRoutes(r, deps)
	return nil
}

func registerExampleRoutes(r *gin.RouterGroup, deps Dependencies) {
	if deps.Example == nil {
		return
	}

	examples := r.Group("/examples")
	examples.GET("", deps.Example.List)
	examples.POST("", deps.Example.Create)
}
