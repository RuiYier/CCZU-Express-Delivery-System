package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yurin-kami/PackChann/controllers"
	"github.com/yurin-kami/PackChann/middlewares"
	"github.com/yurin-kami/PackChann/models"
	"gorm.io/gorm"
)

func ProtectedRoutes(db *gorm.DB, router *gin.Engine, jwtCfg models.JWTConfig) {
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware(db, jwtCfg))
	{
		protected.GET("/getPackDetails/:pack_id", controllers.GetPackDetailsByPackId(db))
		protected.POST("/packCheckIn", controllers.CheckInPack(db))
		protected.POST("/packCheckout", controllers.CheckOutPack(db))
		protected.POST("/mailPack", controllers.MailPack(db))
		protected.POST("/cancelMail", controllers.CancelMailPack(db))
		protected.POST("/updatePackStatus", controllers.UpdatePackStatus(db))
		protected.GET("/allPacks/:user_id", controllers.GetAllPacksByUserId(db))
		protected.POST("/updateUserInfo", controllers.UpdateUserInfoByPhone(db))

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.GET("/users", controllers.GetAllUsers(db))
			admin.GET("/packs", controllers.GetAllPacks(db))
			admin.PUT("/pack", controllers.AdminUpdatePack(db))
			admin.DELETE("/deleteUser", controllers.DeleteUser(db))
		}
	}
}
