package main

import (
	"github.com/Grama-Check/Grama-Check-App/handlers"
	"github.com/Grama-Check/Grama-Check-App/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"*"}
	// config.AllowHeaders = []string{"*"}

	// router.Use(cors.New(config))

	router.Use(static.Serve("/", static.LocalFile("./public", false)))

	authGroup := router.Group("/").Use(middleware.AuthMiddleware())

	authGroup.POST("/gramacheck", handlers.ResponseHandler)
	authGroup.POST("/status", handlers.GetStatus)
	authGroup.GET("/gettoken", handlers.GetToken)
	authGroup.POST("/gramatest", handlers.ResponseHandlerexists)
	router.Run(":9090")
}
