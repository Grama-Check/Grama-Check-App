package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "No token present")
			return
		}

		fields := strings.Fields(token)

		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Single string no prefix")
			return
		}

		if strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Incorrect prefix")
			return
		}

		accessToken := fields[1]

		req, err1 := http.NewRequest("GET", "https://api.asgardeo.io/t/gramacheck/oauth2/userinfo", nil)

		if err1 != nil {
			fmt.Println("error 1")
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err2 := http.DefaultClient.Do(req)

		if err2 != nil {
			fmt.Println("error 2")
		}

		invalidToken := models.InvalidToken{}

		err3 := json.NewDecoder(resp.Body).Decode(&invalidToken)

		if err3 == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Access token validation failed")
			return
		}

		// user := models.AuthorizedUser{}

		// err = json.NewDecoder(resp.Body).Decode(&user)

		// if err != nil {
		// 	fmt.Println("err couldn't read body:", err)
		// 	return
		// }

		fmt.Println("here")

		defer resp.Body.Close()

		// c.Next()
	}
}
