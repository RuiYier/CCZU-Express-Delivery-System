package models

import "time"

type Pack struct {
	PackId       int64     `gorm:"primaryKey;index:idx_pack_id,type:btree" json:"pack_id"`
	UserId       int64     `gorm:"not null;index:idx_packs_user_id,type:btree" json:"user_id"`
	PackStatus   string    `gorm:"type:varchar(20);default:'pending'" json:"pack_status"`
	PickupCode   string    `gorm:"type:varchar(10);index:idx_pickup_code,type:btree" json:"pickup_code"`
	CheckInTime  time.Time `gorm:"autoCreateTime" json:"check_in_time"`
	CheckOutTime time.Time `json:"check_out_time"`
}

type CheckInPak struct {
	PackId    int64 `json:"pack_id" binding:"required"`
	UserId    int64 `json:"user_id" binding:"required"`
	ShelfCode int64 `json:"shelf_code" binding:"required"`
}

type CheckOutPak struct {
	PackId int64 `json:"pack_id" binding:"required"`
	UserId int64 `json:"user_id" binding:"required"`
}

type UpdatePackStatus struct {
	PackId     int64  `json:"pack_id" binding:"required"`
	PackStatus string `gorm:"type:varchar(20);default:'pending'" json:"pack_status"`
}

type AdminUpdatePackInput struct {
	PackId       int64      `json:"pack_id" binding:"required"`
	UserId       *int64     `json:"user_id"`
	PackStatus   *string    `json:"pack_status"`
	PickupCode   *string    `json:"pickup_code"`
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
}

type MailPack struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Recipient       string `json:"recipient" binding:"required"`
	RecivingAddress string `json:"reciving_address" binding:"required"`
	ShipperPhone    string `json:"shipper_phone" binding:"required"`
	RecipientPhone  string `json:"recipient_phone" binding:"required"`
}

func (p *Pack) NewPack(packId int64, userId int64, packStatus, pickupCode string, checkInTime time.Time, checkOutTime time.Time) *Pack {
	return &Pack{
		PackId:       packId,
		UserId:       userId,
		PackStatus:   packStatus,
		PickupCode:   pickupCode,
		CheckInTime:  checkInTime,
		CheckOutTime: checkOutTime,
	}
}
