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

func PoliceCheck(p models.Person, c *gin.Context) {
	reqstr := fmt.Sprintf(`{"nic":"%s","address":"%s","name":"%s"}`, p.NIC, p.Address, p.Name)
	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, PoliceIP, bodyReader)
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

	// resBody, err := ioutil.ReadAll(res.Body)
	policechecked := models.PoliceCheck{}
	err = json.NewDecoder(res.Body).Decode(&policechecked)
	if err != nil {
		fmt.Println("err couldn't read body:", err)
		return
	}
	log.Println("Police: NIC from res", policechecked.NIC, "exists from res", policechecked.Clear)

	log.Println("Police: NIC from res", policechecked.NIC, "exists from res", policechecked.Clear)

	if policechecked.Clear {
		err = queries.UpdatePolice(context.Background(), policechecked.NIC)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		queries.UpdateFailed(context.Background(), db.UpdateFailedParams{Nic: policechecked.NIC, Failed: false})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}

	} else {
		queries.UpdateFailed(context.Background(), db.UpdateFailedParams{Nic: policechecked.NIC, Failed: true})
	}

}
