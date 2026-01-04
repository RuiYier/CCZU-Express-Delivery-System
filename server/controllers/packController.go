package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yurin-kami/PackChann/models"
	"github.com/yurin-kami/PackChann/utils"
	"gorm.io/gorm"
)

func CheckInPack(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var checkInData models.CheckInPak
		if err := c.ShouldBindJSON(&checkInData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var existPack models.Pack
		err := db.WithContext(ctx).Where("pack_id = ? AND pack_status = ?", checkInData.PackId, "pending").First(&existPack).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Pack already checked in"})
			return
		}

		newPack := models.Pack{
			PackId:      checkInData.PackId,
			UserId:      checkInData.UserId,
			PackStatus:  "pending",
			PickupCode:  fmt.Sprint(checkInData.ShelfCode) + "-" + fmt.Sprint(time.Now().UnixNano())[:4], //1-1-1978
			CheckInTime: time.Now(),
		}
		err = db.WithContext(ctx).Create(&newPack).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check in pack"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pack checked in successfully", "pack": newPack})

	}
}

func CheckOutPack(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var checkOutData models.CheckOutPak
		if err := c.ShouldBindJSON(&checkOutData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var pendingPack models.Pack
		err := db.WithContext(ctx).Where("pack_id = ? AND user_id = ? AND pack_status = ?", checkOutData.PackId, checkOutData.UserId, "pending").First(&pendingPack).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No pending pack found for checkout"})
			return
		}

		pendingPack.PackStatus = "checked_out"
		pendingPack.CheckOutTime = time.Now()
		err = db.WithContext(ctx).Save(&pendingPack).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check out pack"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pack checked out successfully", "pack": pendingPack})
	}
}

func GetAllPacksByUserId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var packs []models.Pack
		if err := db.WithContext(ctx).Where("user_id = ?", c.Param("user_id")).Find(&packs).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pack Not found"})
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"packs": packs})
	}
}

func GetPackDetailsByPackId(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var pack models.Pack
		if err := db.WithContext(ctx).Where("pack_id = ?", c.Param("pack_id")).First(&pack).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "pack not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"pack": pack})
	}
}

func CancelMailPack(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cancelMailPack models.CheckOutPak
		if err := c.ShouldBindJSON(&cancelMailPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "input Invalid"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var pack models.Pack
		if err := db.WithContext(ctx).Where("pack_id = ? AND user_id = ? AND pack_status = ?", cancelMailPack.PackId, cancelMailPack.UserId, "in_transit").First(&pack).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "pack not found"})
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		pack.PackStatus = "cancelled"
		if err := db.WithContext(ctx).Save(&pack).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cancelled_mail_pack": pack})
	}
}

func MailPack(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var mailPack models.MailPack
		if err := c.ShouldBindJSON(&mailPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var user models.User
		err := db.WithContext(ctx).Where("phone = ?", mailPack.RecipientPhone).First(&user).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// 生成包裹 ID（雪花算法）
		packId, err := utils.GenerateID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate pack id"})
			return
		}

		newPack := models.Pack{
			PackId:     packId,
			PackStatus: "in_transit",
			UserId:     user.UserId,
		}
		err = db.WithContext(ctx).Create(&newPack).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mail pack"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Mail pack created successfully", "pack": newPack})
	}
}

func UpdatePackStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updatePackStatus models.UpdatePackStatus
		if err := c.ShouldBindJSON(&updatePackStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "input Invalid"})
			return
		}

		ctx, candel := context.WithTimeout(c, 30*time.Second)
		defer candel()

		var pack models.Pack
		if err := db.WithContext(ctx).Where("pack_id = ?", updatePackStatus.PackId).First(&pack).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "pack not found"})
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		pack.PackStatus = updatePackStatus.PackStatus
		if pack.PackStatus == "shipped" {
			pack.CheckOutTime = time.Now()
		}

		if err := db.WithContext(ctx).Save(&pack).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"updatePackStatus": pack})
	}
}

// GetAllPacks 获取所有包裹（管理员权限），支持按状态筛选
func GetAllPacks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := c.Query("status")
		var packs []models.Pack

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		query := db.WithContext(ctx)
		if status != "" {
			query = query.Where("pack_status = ?", status)
		}

		if err := query.Find(&packs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"packs": packs})
	}
}

// AdminUpdatePack 管理员更新包裹信息
func AdminUpdatePack(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.AdminUpdatePackInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c, 30*time.Second)
		defer cancel()

		var pack models.Pack
		if err := db.WithContext(ctx).Where("pack_id = ?", input.PackId).First(&pack).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pack not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		updates := make(map[string]interface{})
		if input.UserId != nil {
			updates["user_id"] = *input.UserId
		}
		if input.PackStatus != nil {
			updates["pack_status"] = *input.PackStatus
		}
		if input.PickupCode != nil {
			updates["pickup_code"] = *input.PickupCode
		}
		if input.CheckInTime != nil {
			updates["check_in_time"] = *input.CheckInTime
		}
		if input.CheckOutTime != nil {
			updates["check_out_time"] = *input.CheckOutTime
		}

		if len(updates) > 0 {
			if err := db.WithContext(ctx).Model(&pack).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pack"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"pack": pack})
	}
}
