package main

import (
	"bytes"
	"net/http"

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

	resjson := IdentityCheck(c)

	c.JSON(
		http.StatusOK,
		resjson,
	)
}
