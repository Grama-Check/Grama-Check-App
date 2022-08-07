package handlers

import (
	"context"
	"database/sql"
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
	//log.Println(config.SendGridKey, ":", config.DBSource)
	conn, err := sql.Open(config.DBDriver, config.DBSource)

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
