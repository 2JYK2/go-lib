package mgrModel

import (
	"time"
)

// ConfigDeviceType Server device model configuration table
type ConfigDeviceType struct {
	ID               int       `gorm:"primaryKey;column:id;type:int;not null" json:"-"`                                        // Auto increment primary key
	Type             string    `gorm:"column:type;type:varchar(300);not null" json:"type"`                                     // Container automatically reports the model string
	Term             string    `gorm:"index:idx_term;column:term;type:varchar(30);not null;default:''" json:"term"`            // Model abbreviation
	Power            int       `gorm:"column:power;type:int;not null;default:0" json:"power"`                                  // Hardware model corresponding rendering computing power label value
	PowerVideoEngine int       `gorm:"column:power_video_engine;type:int;not null;default:0" json:"power_video_engine"`        // Reserved field, hardware model corresponding to video engine capability label value
	SupportVp8Decode bool      `gorm:"column:support_vp8_decode;type:tinyint(1);not null;default:0" json:"support_vp8_decode"` // Whether the GPU supports VP8 decoding, 1: supported, 0: not supported
	Note             string    `gorm:"column:note;type:varchar(200);not null" json:"note"`                                     // Remark explanation
	CreateAt         time.Time `gorm:"column:create_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"create_at"`    // Data record creation time
	UpdateAt         time.Time `gorm:"column:update_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"update_at"`    // Data update time
}

// TableName get sql table name.
func (m *ConfigDeviceType) TableName() string {
	return "config_device_type"
}
