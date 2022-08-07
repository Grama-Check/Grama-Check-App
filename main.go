package main

import (
	"github.com/Grama-Check/Grama-Check-App/handlers"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(static.Serve("/", static.LocalFile("./public", false)))

	router.POST("/gramacheck", handlers.ResponseHandler)
	router.POST("/status", handlers.GetStatus)
	router.GET("/gettoken", handlers.GetToken)
	router.POST("/create", handlers.CreateUser)
	router.POST("/gramatest", handlers.ResponseHandlerexists)
	router.Run(":9090")
}
