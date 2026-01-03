package models

import (
	"time"
)

type User struct {
	UserId       int64     `gorm:"primaryKey;autoIncrement" json:"user_id"`
	UserName     string    `gorm:"type:varchar(100);not null" json:"user_name"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	StudentId    string    `gorm:"type:varchar(50);unique;not null" json:"student_id"`
	Phone        string    `gorm:"type:varchar(20);unique;not null" json:"phone"`
	Address      string    `gorm:"type:varchar(255)" json:"address"`
	Role         string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	RegisterTime time.Time `gorm:"autoCreateTime" json:"register_time"`
}

type UserLogin struct {
	StudentId string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (u *User) NewUser(userId int64, userName, passwordHash, studentId, phone string, address string, registerTime time.Time) *User {
	return &User{
		UserId:       userId,
		UserName:     userName,
		PasswordHash: passwordHash,
		StudentId:    studentId,
		Phone:        phone,
		Address:      address,
		RegisterTime: registerTime,
	}
}
