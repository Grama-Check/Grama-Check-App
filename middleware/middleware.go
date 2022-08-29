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
		// Checking if a token is sent in the Authorization Header
		token := c.GetHeader("Authorization")
		
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "No token present")
			return
		}
		
		
		fields := strings.Fields(token)
		
		// Checking if a token has prefix bearer
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Single string no prefix")
			return
		}
		
		// Checking if the prefix is equal to "brearer"
		if strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Incorrect prefix")
			return
		}

		accessToken := fields[1]
		
		// Getting user details from asgardeo
		req, err1 := http.NewRequest("GET", "https://api.asgardeo.io/t/gramacheck/oauth2/userinfo", nil)
		if err1 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err1.Error())
			return
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err2 := http.DefaultClient.Do(req)
		if err2 != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err2.Error())
			return
		}
		
		// If the expected response is not recieved an issue with access token
		user := models.AuthorizedUser{}

		err4 := json.NewDecoder(resp.Body).Decode(&user)
		if err4 != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err4.Error())
			return
		}
		
		// If the expected response is 
		if (models.AuthorizedUser{} == user) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Access token validation failed")
			return
		}

		formType := c.GetHeader("Form")
		
		// Checking if the request form was sent by submitting the apply form or check form
		if formType == "apply" {
			model := models.Person{}

			if err := c.ShouldBindBodyWith(&model, binding.JSON); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			formNIC := strings.ReplaceAll(model.NIC, " ", "")
			userNIC := strings.ReplaceAll(user.NIC, " ", "")

			if !strings.EqualFold(formNIC, userNIC) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "NIC mismatch")
				return
			}
		} else if formType == "check" {
			model := models.StatusCheck{}

			if err := c.ShouldBindBodyWith(&model, binding.JSON); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			formNIC := strings.ReplaceAll(model.NIC, " ", "")
			userNIC := strings.ReplaceAll(user.NIC, " ", "")

			if !strings.EqualFold(formNIC, userNIC) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "NIC mismatch")
				return
			}
		}

		c.Next()
	}
}
