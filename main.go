package main

import (
	"github.com/Grama-Check/Grama-Check-App/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	auth.GenerateKeys()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", Index)
	router.Static("/public/well-known", "./public/well-known")
	//router.GET("/signup", Signup)
	//router.GET("/login", Login)

	//authGroup := router.Group("/").Use()
	//router.GET("/home", Home)
	router.POST("/gramacheck", responseHandler)
	// -> Identity Check , id, address , UID
	// <- Failed/Passed , UID , Pass/Fail
	router.Run(":9090")
}
