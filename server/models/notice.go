package models

import "time"

type Notice struct {
	NoticeId  int64     `gorm:"primaryKey;autoIncrement" json:"notice_id"`
	CompanyId int32     `gorm:"not null" json:"company_id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	SendTime  time.Time `gorm:"autoCreateTime" json:"send_time"`
}

func (n *Notice) NewNotice(noticeId int64, companyId int32, title, content string, sendTime time.Time) *Notice {
	return &Notice{
		NoticeId:  noticeId,
		CompanyId: companyId,
		Title:     title,
		Content:   content,
		SendTime:  sendTime,
	}
}
