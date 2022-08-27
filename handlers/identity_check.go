package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Grama-Check/Grama-Check-App/auth"
	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/Grama-Check/Grama-Check-App/util"
	"github.com/gin-gonic/gin"
)

func IdentityCheck(p models.Person, c *gin.Context) {
	// preparing the request to send to identity check mico service
	reqstr := fmt.Sprintf(`{"nic":"%s"}`, p.NIC)
	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, IdentityIP, bodyReader)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// generating the token
	token, err := auth.GenerateToken()
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, "Couldn't generate token")
		return
	}

	// setting the Authorization header
	authHeader := fmt.Sprintf("Bearer %v", token)

	req.Header.Add("Authorization", authHeader)

	// making the request to identity check micro service
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the expected result is recieved from identity check
	idchecked := models.IDChecked{}

	err = json.NewDecoder(res.Body).Decode(&idchecked)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("Identity: NIC from res", idchecked.NIC, "exists from res", idchecked.Exists)

	// updating the checks table based on the clear status
	if idchecked.Exists {
		err = queries.UpdateIdentityCheck(context.Background(), idchecked.NIC)
		if err != nil {
			util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		// calling address check method
		Addresscheck(p, c)
	} else {
		util.SendIssue(p, "Identity")
		queries.UpdateFailed(context.Background(), idchecked.NIC)
	}
}
