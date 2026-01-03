package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yurin-kami/PackChann/controllers"
	"github.com/yurin-kami/PackChann/models"
	"gorm.io/gorm"
)

func UnprotectedRoutes(db *gorm.DB, router *gin.Engine, jwtCfg models.JWTConfig) {
	// Define your unprotected routes here
	router.POST("/register", controllers.RegisterUser(db, jwtCfg))
	router.POST("/login", controllers.LoginUser(db, jwtCfg))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}
