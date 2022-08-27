package main

import (
	"main/db"
	"main/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	res := db.PrepareDB()
	print(res)
	app.GET("/", routers.Root)
	farmer := app.Group("/api/v1")
	farmer.GET("/loginservice/:contact_num", routers.LoginService)
	farmer.GET("/farmer/fullprofile/:contact_num", routers.FullProfile)
	farmer.POST("/loginservice/newuser", routers.CreateNewUser)
	app.Run()
}
