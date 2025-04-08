package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/login", controllers.Login())
	incomingRoutes.POST("/signup", controllers.SignUp())
}
