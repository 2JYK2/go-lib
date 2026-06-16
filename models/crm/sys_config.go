package crmModel

import (
	"time"
)

type SysConfig struct {
	ID        uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ConfName  string    `gorm:"column:conf_name;NOT NULL"`
	ConfValue string    `gorm:"column:conf_value;NOT NULL"`
	Module    string    `gorm:"column:module;NOT NULL"`
	CreateAt  time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdateAt  time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *SysConfig) TableName() string {
	return "sys_config"
}
