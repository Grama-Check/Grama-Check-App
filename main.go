package main

import (
	"github.com/Grama-Check/Grama-Check-App/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", handlers.Index)
	router.Static("/public/well-known", "./public/well-known")
	//router.GET("/signup", Signup)
	//router.GET("/login", Login)

	//authGroup := router.Group("/").Use()
	//router.GET("/home", Home)
	router.POST("/gramacheck", handlers.ResponseHandler)
	router.POST("/status", handlers.GetStatus)
	router.GET("/gettoken", handlers.GetToken)
	router.POST("/create", handlers.CreateUser)
	// -> Identity Check , id, address , UID
	// <- Failed/Passed , UID , Pass/Fail
	router.Run(":9090")
}
