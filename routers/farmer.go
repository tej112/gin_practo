package routers

import (
	"fmt"
	"main/db"
	"main/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "login-farmer service is up and running",
	})
}
func LoginService(c *gin.Context) {
	num := c.Param("contact_num")
	c1 := make(chan string)
	go func() {
		user := db.Exists(num)
		if user {
			c1 <- "existing_user"
		} else {
			c1 <- "new_user"
		}
	}()

	c2 := make(chan map[string]string)

	go func() {
		c2 <- db.HashOfFarmer(num)

	}()
	hash := <-c2
	usr := <-c1

	c.JSON(http.StatusOK, gin.H{
		"flag":     1,
		"response": usr,
		"info":     hash,
	})
}

func FullProfile(c *gin.Context) {
	pk := db.ValueForField(c.Param("contact_num"), "pk")
	farmer := db.GetFarmer(pk)

	if farmer != false {
		c.JSON(http.StatusOK, farmer)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "User NotFound",
		})
	}

}

func CreateNewUser(c *gin.Context) {

	newFarmer := models.Farmer{}
	err := c.BindJSON(&newFarmer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"required": strings.Split(strings.Split(err.Error(), " ")[1], "'")[1],
			"error":    err.Error()})
		return
	}
	newFarmer.Pk = ulid.Make().String()
	newFarmer.Created_at = time.Now().Format("2006-01-02 15:04:05")
	newFarmer.Last_updated_at = time.Now().Format("2006-01-02 15:04:05")
	if newFarmer.Profile_pic == "" {
		newFarmer.Profile_complete_perc = 80
	} else {
		newFarmer.Profile_complete_perc = 100
	}
	user := db.Exists(newFarmer.Contact_num)
	if user {
		c.JSON(http.StatusForbidden, gin.H{
			"detail": fmt.Sprintf("user for mobile number %v already exists", newFarmer.Contact_num),
		})
		return
	}

	go db.CreateNewFarmer(newFarmer)

	c.JSON(http.StatusCreated, newFarmer)
}
