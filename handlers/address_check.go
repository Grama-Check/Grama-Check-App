package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Grama-Check/Grama-Check-App/auth"
	db "github.com/Grama-Check/Grama-Check-App/db/sqlc"
	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/gin-gonic/gin"
)

func Addresscheck(p models.Person, c *gin.Context) {
	reqstr := fmt.Sprintf(`{"nic":"%s","address":"%s"}`, p.NIC, p.Address)

	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, addresscheckIP, bodyReader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return

	}

	token, err := auth.GenerateToken()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Couldn't generate token")
		return

	}

	authHeader := fmt.Sprintf("Bearer %v", token)

	req.Header.Add("Authorization", authHeader)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Couldn't perform request:", err)
	}

	addresschecked := models.AddressChecked{}
	err = json.NewDecoder(res.Body).Decode(&addresschecked)
	if err != nil {
		fmt.Println("err couldn't read body:", err)
		return
	}
	log.Println("Address: NIC from res", addresschecked.NIC, "exists from res", addresschecked.Exists)

	log.Println("Address: NIC from res", addresschecked.NIC, "exists from res", addresschecked.Exists)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}

	if addresschecked.Exists {
		err = queries.UpdateAddress(context.Background(), addresschecked.NIC)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		PoliceCheck(p, c)
	} else {
		queries.UpdateFailed(context.Background(), db.UpdateFailedParams{Nic: addresschecked.NIC, Failed: true})
	}

}
