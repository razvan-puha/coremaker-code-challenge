package main

import (
	"demo/coremaker/docs"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const baseUriPath = "/api/v1"

func main() {
	InitDB()

	docs.SwaggerInfo.Title = "Coremaker API Swagger"
	docs.SwaggerInfo.Description = "This is a sample auth server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := gin.Default()
	// For healthcheck
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST(fmt.Sprintf("%s/auth/register", baseUriPath), RegisterUser)
	r.POST(fmt.Sprintf("%s/auth/login", baseUriPath), LoginUser)
	r.GET(fmt.Sprintf("%s/auth/currentUser", baseUriPath), VerifyToken, GetCurrentUserDetails)

	r.Run()
}
