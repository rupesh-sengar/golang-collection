package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rupesh-sengar/golang-collection/auth/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/login", controllers.LoginHandler)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		api.POST("/signup", controllers.SignupHandler)
		api.POST("/approve-user", controllers.AuthApprovalHandler)
	}
}


