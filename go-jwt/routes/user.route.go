package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/controllers"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:id", controllers.GetUser())
}
