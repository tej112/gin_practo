package routers

import (
	"fmt"
	"log"
	"main/db"
	"main/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"github.com/oklog/ulid/v2"
)

func checkErrRedis(c *gin.Context, err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DataBase Connection Lost. Try Again later"})
			return true
		} else if strings.Contains(err.Error(), "timeout") {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Connection Timed Out",
			})
			return true
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err, err.Error())
		return true
	}
	return false
}

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "login-farmer service is up and running",
	})
}

func LoginService(c *gin.Context) {
	num := c.Param("contact_num") //get the mobile number from the endpoint and store it in "num" variable

	// create a channel c1 and run if mobile exists concurrently and store the result in the c1
	c1 := make(chan *redis.BoolCmd)
	go func() {
		user := db.Exists(num)
		c1 <- user
	}()

	// create a channel c2 and run func to get farmer's hash form db and strore the result in c2
	c2 := make(chan *redis.StringStringMapCmd)
	go func() {
		c2 <- db.HashOfFarmer(num)

	}()

	//extract the info from c1 and store it in "exists"
	exists := <-c1
	// extract the response/result and error, and check if any connection error occured
	res, err := exists.Result()
	if checkErrRedis(c, err) {
		return
	}

	//if only user exists,then get the hash of farmer and check for error while getting hash of farmer
	if res {
		hash := <-c2
		hashres, hasherr := hash.Result()
		if checkErrRedis(c, hasherr) {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"flag":     1,
			"response": res,
			"info":     hashres,
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"flag":     1,
		"response": res,
		"info":     gin.H{},
	})
}

func FullProfile(c *gin.Context) {
	pk := db.ValueForField(c.Param("contact_num"), "pk")

	res, err := pk.Result()
	if checkErrRedis(c, err) {
		return
	}

	farmer := db.GetFarmer(res)

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
	err := c.ShouldBindJSON(&newFarmer)
	if err != nil {
		if err.Error() == "EOF" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "No data is provided",
			})
			return
		}
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
	res, err := user.Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if res {
		c.JSON(http.StatusForbidden, gin.H{
			"detail": fmt.Sprintf("user for mobile number %v already exists", newFarmer.Contact_num),
		})
		return
	}

	go db.CreateNewFarmer(newFarmer)

	c.JSON(http.StatusCreated, newFarmer)
}
