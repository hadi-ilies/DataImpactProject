package main

import (
	"DataImpactProject/server/app/controllers/usercontrollers"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	api := r.Group("/")
	api.POST("/add/users", usercontrollers.Create)
	api.POST("/login", usercontrollers.Login)
	api.DELETE("/delete/user/:id", usercontrollers.Delete)
	api.GET("/users/list", usercontrollers.GetAllUsers)
	api.GET("/users/:id", usercontrollers.GetUserByID)
	api.PUT("/user/:id", usercontrollers.Update)
}
