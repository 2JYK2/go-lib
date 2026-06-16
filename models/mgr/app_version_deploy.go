package mgrModel

import (
	"time"
)

type AppVersionDeploy struct {
	AppID         int64     `gorm:"column:app_id;NOT NULL"`
	TargetVersion int64     `gorm:"column:target_version;NOT NULL"`
	RegionCode    string    `gorm:"column:region_code;NOT NULL"`
	ArsID         int64     `gorm:"column:ars_id;NOT NULL"`
	Status        string    `gorm:"column:status;NOT NULL"`
	CreateAt      time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdateAt      time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *AppVersionDeploy) TableName() string {
	return "app_version_deploy"
}
