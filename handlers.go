package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Grama-Check/Grama-Check-App/auth"
	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/gin-gonic/gin"
)

const (
	IdentityIP = "http://localhost:8080"
)

func Index(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "testing",
		},
	)
}

func IdentityCheck(c *gin.Context) models.IDChecked {
	jsonBody := []byte(`{
		"uid":123,
		"id":"izwcvpjmik"
	}`)

	bodyReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest(http.MethodPost, IdentityIP, bodyReader)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Failed request to Identity check")

	}
	token, err := auth.GenerateToken()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Failed to generate token")

	}
	authHeader := fmt.Sprintf("Bearer %v", token)
	c.Header("authorization", authHeader)

	resjson := models.IDChecked{}
	c.BindJSON(res.Body)

	return resjson
}

func responseHandler(c *gin.Context) {

	person := models.Person{}

	err := c.BindJSON(&person)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Couldn't parse request body")
		return
	}

	// Add to database if they don't exist

	// check identity if exists return and add to database

	// check address if exists return and add to db

	// check for criminal record, add to database

	// send completion

	resjson := IdentityCheck(c)

	c.JSON(
		http.StatusOK,
		resjson,
	)
}

func getStatus(c *gin.Context) {

}
