package routes

import (
	"github.com/gin-gonic/gin"
)

// @tag.name Test
// @Summary ping
// @Schemes
// @Description do ping
// @Tags Test
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /ping [get]
func SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

	}

	return r
}
