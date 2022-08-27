package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// checking if a token is present in the Authorization Header
		token := c.GetHeader("Authorization")

		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "No token present")
			return
		}

		fields := strings.Fields(token)

		// checking if a token has prefix bearer
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Single string no prefix")
			return
		}

		// checking if the prefix is equal to "brearer"
		if strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Incorrect prefix")
			return
		}

		accessToken := fields[1]

		// getting user details from asgardeo
		// configuring the request
		req, err1 := http.NewRequest("GET", "https://api.asgardeo.io/t/gramacheck/oauth2/userinfo", nil)
		if err1 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err1.Error())
			return
		}

		// setting the authorization header in the request
		req.Header.Set("Authorization", "Bearer "+accessToken)

		// making the request to asgardeo
		resp, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err2.Error())
			return
		}

		// struct to store information receiving from Asgardeo
		user := models.AuthorizedUser{}

		err4 := json.NewDecoder(resp.Body).Decode(&user)
		if err4 != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err4.Error())
			return
		}

		// if the expected response is not recieved an issue with access token
		if (models.AuthorizedUser{} == user) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Access token validation failed")
			return
		}

		formType := c.GetHeader("Form")

		// checking if the request form was sent by submitting the apply form or check form
		if formType == "apply" {
			// struct to store infotmation recived upon apply form submission
			model := models.Person{}

			if err := c.ShouldBindBodyWith(&model, binding.JSON); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			// removing leading trailing whitespaces from the NICs
			formNIC := strings.ReplaceAll(model.NIC, " ", "")
			userNIC := strings.ReplaceAll(user.NIC, " ", "")

			// comparing the NICs received from Asgardeo and the form
			if !strings.EqualFold(formNIC, userNIC) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "NIC mismatch")
				return
			}
		} else if formType == "check" {
			// struct to store infotmation recived upon check form submission
			model := models.StatusCheck{}

			if err := c.ShouldBindBodyWith(&model, binding.JSON); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			// removing leading trailing whitespaces from the NICs
			formNIC := strings.ReplaceAll(model.NIC, " ", "")
			userNIC := strings.ReplaceAll(user.NIC, " ", "")

			// comparing the NICs received from Asgardeo and the form
			if !strings.EqualFold(formNIC, userNIC) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "NIC mismatch")
				return
			}
		}

		c.Next()
	}
}
