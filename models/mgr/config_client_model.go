package mgrModel

import "time"

type ConfigClientModel struct {
	ID                      int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	ClientModel             string    `gorm:"column:client_model;NOT NULL" json:"client_model"`                            // Terminal device model, terminal device model identification, composed of four elements of terminal device, basic format:
	IsAllowJoin             int       `gorm:"column:is_allow_join" json:"is_allow_join"`                                   // Allow access to the system
	SupportVideoTypes       string    `gorm:"column:support_video_types" json:"support_video_types"`                       // Video types H.265, H.264, vp8
	SupportAudioTypes       string    `gorm:"column:support_audio_types" json:"support_audio_types"`                       // Audio type aac, opus
	Comment                 string    `gorm:"column:comment" json:"comment"`                                               // memo
	HardwareType            string    `gorm:"column:hardware_type" json:"hardware_type"`                                   // Hardware type, chip: chip, machine: complete machine
	Manufacturer            string    `gorm:"column:manufacturer" json:"manufacturer"`                                     // Equipment model manufacturer
	EnableEnhanceResolution int       `gorm:"column:enable_enhance_resolution;default:0" json:"enable_enhance_resolution"` // Whether to enable resolution superresolution capability
	CreateAt                time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"`        // Data record creation time
	UpdateAt                time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"`        // Data update time
	MaxResoLevel            string    `gorm:"column:max_reso_level" json:"max_reso_level"`                                 // Maximum clarity level supported
}

func (m *ConfigClientModel) TableName() string {
	return "config_client_model"
}
