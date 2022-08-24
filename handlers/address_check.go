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

func Addresscheck(p models.Person, c *gin.Context) {
	reqstr := fmt.Sprintf(`{"nic":"%s","address":"%s"}`, p.NIC, p.Address)

	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, addresscheckIP, bodyReader)
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

	addresschecked := models.AddressChecked{}

	err = json.NewDecoder(res.Body).Decode(&addresschecked)
	if err != nil {
		util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("Address: NIC from res", addresschecked.NIC, "exists from res", addresschecked.Exists)

	// log.Println("Address: NIC from res", addresschecked.NIC, "exists from res", addresschecked.Exists)

	if addresschecked.Exists {
		err = queries.UpdateAddressCheck(context.Background(), addresschecked.NIC)
		if err != nil {
			util.SendError(http.StatusInternalServerError, p.NIC+" "+err.Error())

			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		PoliceCheck(p, c)
	} else {
		util.SendIssue(p, "Address")
		queries.UpdateFailed(context.Background(), addresschecked.NIC)
	}
}
