package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/Grama-Check/Grama-Check-App/auth"
	db "github.com/Grama-Check/Grama-Check-App/db/sqlc"
	"github.com/Grama-Check/Grama-Check-App/models"
	"github.com/Grama-Check/Grama-Check-App/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/lib/pq"
)

const (
	IdentityIP     = "https://identity-check-service.mangobeach-b9b75009.westus.azurecontainerapps.io"
	addresscheckIP = "https://address-check-service.mangobeach-b9b75009.westus.azurecontainerapps.io"
	PoliceIP       = "https://police-check-service.mangobeach-b9b75009.westus.azurecontainerapps.io"
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
		util.SendError(http.StatusInternalServerError, err2.Error())

		// log.Println(http.StatusInternalServerError, err.Error(), "", err2.Error())
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

	if err := c.ShouldBindBodyWith(&person, binding.JSON); err != nil {
		util.SendError(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	args := db.CreateCheckParams{
		Nic:          person.NIC,
		Address:      person.Address,
		Email:        person.Email,
		Name:         person.Name,
		Idcheck:      false,
		Addresscheck: false,
		Policecheck:  false,
		Failed:       false,
	}

	_, err := queries.CreateCheck(context.Background(), args)
	if err != nil {
		duplicateError := "pq: duplicate key value violates unique constraint \"checks_pkey\""

		if strings.EqualFold(err.Error(), duplicateError) {
			err = queries.DeleteCheck(context.Background(), args.Nic)
			if err != nil {
				util.SendError(http.StatusInternalServerError, person.NIC+" "+err.Error())

				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}

			_, err = queries.CreateCheck(context.Background(), args)
			if err != nil {
				util.SendError(http.StatusInternalServerError, person.NIC+" "+err.Error())

				c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			util.SendError(http.StatusInternalServerError, person.NIC+" "+err.Error())

			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	go IdentityCheck(person, c)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}

func GetStatus(c *gin.Context) {
	nic := models.StatusCheck{}

	if err := c.ShouldBindBodyWith(&nic, binding.JSON); err != nil {
		util.SendError(http.StatusBadRequest, nic.NIC+" "+err.Error())

		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	check, err := queries.GetCheck(context.Background(), nic.NIC)
	if err != nil {
		util.SendError(http.StatusInternalServerError, nic.NIC+" "+err.Error())

		missingError := "sql: no rows in result set"

		if strings.EqualFold(err.Error(), missingError) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "Please send a grama check request first")
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"exists":       true,
			"name":         check.Name,
			"nic":          check.Nic,
			"idcheck":      check.Idcheck,
			"policecheck":  check.Policecheck,
			"addresscheck": check.Addresscheck,
		},
	)
}

func GetToken(c *gin.Context) {
	token, err := auth.GenerateToken()

	if err != nil {
		util.SendError(http.StatusInternalServerError, err.Error())

		c.AbortWithStatusJSON(http.StatusInternalServerError, "Couldn't generate token: "+err.Error())
		return
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
		util.SendError(http.StatusBadRequest, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, `{"error";"Couldnt parse request to json}"`)
		return
	}

	go IdentityCheck(person, c)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "ok",
		},
	)
}
