package crmModel

import "time"

type SnWhite struct {
	ID       int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Tid      int64     `gorm:"column:tid;NOT NULL"`
	Sn       string    `gorm:"column:sn;NOT NULL"`
	CreateAt time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdateAt time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *SnWhite) TableName() string {
	return "sn_white"
}
