package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yurin-kami/PackChann/models"
	"github.com/yurin-kami/PackChann/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func passwordHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// 验证密码
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// 保存 Token 到数据库
func saveTokenToDB(ctx context.Context, db *gorm.DB, userId int64, accessToken, refreshToken string, expirationHours int16) error {
	expiresAt := time.Now().Add(time.Duration(expirationHours) * time.Hour)
	userToken := models.UserToken{
		UserId:       userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	return db.WithContext(ctx).Create(&userToken).Error
}

type RegisterInput struct {
	UserName  string `json:"user_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	StudentId string `json:"student_id" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address"`
	Role      string `json:"role"` // Optional, default to 'user'
}

func RegisterUser(db *gorm.DB, jwtCfg models.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		role := "user"
		if input.Role == "admin" {
			role = "admin"
		}

		newUser := models.User{
			UserName:     input.UserName,
			PasswordHash: passwordHash(input.Password),
			StudentId:    input.StudentId,
			Phone:        input.Phone,
			Address:      input.Address,
			Role:         role,
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var existingUser models.User
		err := db.WithContext(ctx).Where("student_id = ? OR phone = ?", newUser.StudentId, newUser.Phone).First(&existingUser).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 生成用户 ID（雪花算法）并设置
		if id, err := utils.GenerateID(); err == nil {
			newUser.UserId = id
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate user id"})
			return
		}

		if err := db.WithContext(ctx).Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 生成 Token
		accessToken, refreshToken, err := utils.GenerateAllToken(newUser, jwtCfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// 存储 Token
		if err := saveTokenToDB(ctx, db, newUser.UserId, accessToken, refreshToken, jwtCfg.ExpirationHours); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":          newUser,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func LoginUser(db *gorm.DB, jwtCfg models.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginInput models.UserLogin
		if err := c.ShouldBindJSON(&loginInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()
		var user models.User
		err := db.WithContext(ctx).Where("student_id = ?", loginInput.StudentId).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// 使用 bcrypt 验证密码
		if !checkPassword(user.PasswordHash, loginInput.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// 生成 Token
		accessToken, refreshToken, err := utils.GenerateAllToken(user, jwtCfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// 存储 Token
		if err := saveTokenToDB(ctx, db, user.UserId, accessToken, refreshToken, jwtCfg.ExpirationHours); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": models.User{
				UserId:    user.UserId,
				UserName:  user.UserName,
				Phone:     user.Phone,
				Address:   user.Address,
				StudentId: user.StudentId,
				Role:      user.Role,
			},
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func UpdateUserInfoByPhone(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateUser models.User
		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "input Invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var existingUser models.User
		if err := db.WithContext(ctx).Where("user_id = ?", updateUser.UserId).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if err := db.WithContext(ctx).Model(&existingUser).Updates(&updateUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"update_user": existingUser})
	}
}

// GetAllUsers 获取所有用户（管理员权限）
func GetAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		if err := db.WithContext(ctx).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		if err := db.WithContext(ctx).Where("user_id = ?", c.Param("user_id")).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if err := db.WithContext(ctx).Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "delete user complete"})
	}
}
