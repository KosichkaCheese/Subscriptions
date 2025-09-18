package routes

import (
	"subscriptions/handlers"

	"github.com/gin-gonic/gin"
)

// @Summary ping
// @Schemes
// @Description do ping
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping [get]
func SetupRouter(serviceHandler *handlers.ServiceHandler, subscriptionHandler *handlers.SubscriptionHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		api.POST("/services", serviceHandler.Create)
		api.GET("/services", serviceHandler.GetAll)
		api.DELETE("/services/:id", serviceHandler.Delete)

		api.POST("/subs", subscriptionHandler.Create)
		api.GET("/subs", subscriptionHandler.GetAll)
		api.PUT("/subs/:id", subscriptionHandler.Update)
		api.DELETE("/subs/:id", subscriptionHandler.Delete)
		api.GET("/subs/:id", subscriptionHandler.GetById)
		api.GET("/subs/sum", subscriptionHandler.SumByFilters)

	}

	return r
}
