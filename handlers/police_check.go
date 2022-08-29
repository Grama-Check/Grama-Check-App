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

func PoliceCheck(p models.Person, c *gin.Context) {
	reqstr := fmt.Sprintf(`{"nic":"%s","address":"%s","name":"%s"}`, p.NIC, p.Address, p.Name)
	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, PoliceIP, bodyReader)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, "Couldn't generate token")
		return
	}

	authHeader := fmt.Sprintf("Bearer %v", token)

	req.Header.Add("Authorization", authHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// resBody, err := ioutil.ReadAll(res.Body)
	policechecked := models.PoliceCheck{}

	err = json.NewDecoder(res.Body).Decode(&policechecked)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("Police: NIC from res", policechecked.NIC, "exists from res", policechecked.Clear)

	// log.Println("Police: NIC from res", policechecked.NIC, "exists from res", policechecked.Clear)

	if policechecked.Clear {
		err = queries.UpdatePoliceCheck(context.Background(), policechecked.NIC)
		SendStatus(p.NIC)
		if err != nil {
			util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	} else {

		util.SendIssue(p, "Police")
		queries.UpdateFailed(context.Background(), policechecked.NIC)
		SendStatus(p.NIC)
	}
}
