package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Grama-Check/Grama-Check-App/auth"
	db "github.com/Grama-Check/Grama-Check-App/db/sqlc"
	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/Grama-Check/Grama-Check-App/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	IdentityIP     = "http://20.245.188.111:8080"
	addresscheckIP = "http://20.66.32.88:7070"
	PoliceIP       = "http://20.245.209.212:6060"
	dbDriver       = "postgres"
	dbSource       = "postgres://jhivan:25May2001@grama-check-db.postgres.database.azure.com/postgres?sslmode=require"
)

var queries *db.Queries
var config util.Config

func init() {
	var err error
	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	log.Println(config.SendGridKey, ":", config.DBSource)
	conn, err := sql.Open(dbDriver, dbSource)

	err2 := conn.Ping()

	if err != nil || err2 != nil {
		log.Println(http.StatusInternalServerError, `{"error": "couldn't connect to database"`)
		return
	}

	queries = db.New(conn)
}

func Index(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "testing",
		},
	)
}

func ResponseHandler(c *gin.Context) {

	person := models.Person{}

	err := c.BindJSON(&person)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, `{"error";"Couldnt parse request to json}"`)

	}

	args := db.CreateUserParams{
		Nic:          person.NIC,
		Address:      person.Address,
		Email:        person.Email,
		Name:         person.Name,
		Idcheck:      false,
		Addresscheck: false,
		Policecheck:  false,
		Failed:       false,
	}

	_, err = queries.CreateUser(context.Background(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintln("Couldn't add to db: ", err))
		return
	}

	IdentityCheck(person, c)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}

func IdentityCheck(p models.Person, c *gin.Context) {
	reqstr := fmt.Sprintf(`{"nic":"%s"}`, p.NIC)
	jsonBody := []byte(reqstr)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, IdentityIP, bodyReader)
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
	idchecked := models.IDChecked{}
	err = json.NewDecoder(res.Body).Decode(&idchecked)
	if err != nil {
		fmt.Println("err couldn't read body:", err)
		return
	}
	log.Println("Identity: NIC from res", idchecked.NIC, "exists from res", idchecked.Exists)

	log.Println("Identity: NIC from res", idchecked.NIC, "exists from res", idchecked.Exists)

	if idchecked.Exists {
		err = queries.UpdateID(context.Background(), idchecked.NIC)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}

		Addresscheck(p, c)
		return
	} else {
		queries.UpdateFailed(context.Background(), db.UpdateFailedParams{Nic: idchecked.NIC, Failed: true})
	}

}
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

func GetStatus(c *gin.Context) {
	nic := models.StatusCheck{}
	ctx := context.Background()
	c.BindJSON(&nic)

	person, err := queries.GetUser(ctx, nic.NIC)
	if err == nil {
		SendStatus(person)
	}

	if err == nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"exists":       true,
				"name":         person.Name,
				"nic":          person.Nic,
				"failed":       person.Failed,
				"idcheck":      person.Idcheck,
				"policecheck":  person.Policecheck,
				"addresscheck": person.Addresscheck,
			},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"exists": false,
			},
		)
	}

}
func CreateUser(c *gin.Context) {
	person := models.Person{}

	err := c.BindJSON(&person)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, `{"error";"Couldnt parse request to json}"`)

	}

	args := db.CreateUserParams{
		Nic:          person.NIC,
		Address:      person.Address,
		Email:        person.Email,
		Name:         person.Name,
		Idcheck:      false,
		Addresscheck: false,
		Policecheck:  false,
		Failed:       false,
	}

	_, err = queries.CreateUser(context.Background(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, fmt.Sprintln("Couldn't add to db: ", err))
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}

func GetToken(c *gin.Context) {
	token, err := auth.GenerateToken()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Couldn't generate token: "+err.Error())
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			token: token,
		},
	)
}

func ResponseHandlerexists(c *gin.Context) {

	person := models.Person{}

	err := c.BindJSON(&person)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, `{"error";"Couldnt parse request to json}"`)

	}

	IdentityCheck(person, c)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}
