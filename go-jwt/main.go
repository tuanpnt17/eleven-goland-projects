package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanpnt17/eleven-golang-projects/go-jwt/routes"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8686"
	}

	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	_ = router.Run(":" + port)

}
