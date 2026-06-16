package mgrModel

import "time"

type AppVersion struct {
	AppID           int64     `gorm:"column:app_id;primary_key"`                           // Unique application identifier
	Version         int64     `gorm:"column:version;NOT NULL"`                             // Application version
	VersionAlias    string    `gorm:"column:version_alias"`                                // Version alias, game's own version identifier
	SrcDownloadUrl  string    `gorm:"column:src_download_url;NOT NULL"`                    // Source file download URL
	SrcMd5          string    `gorm:"column:src_md5;NOT NULL"`                             // Source file MD5 value
	SrcSize         int       `gorm:"column:src_size;default:0;NOT NULL"`                  // Original resource package size, in MB
	YxConfigContent string    `gorm:"column:yx_config_content;NOT NULL"`                   // Application description file XML content
	Status          string    `gorm:"column:status;NOT NULL"`                              // create: Initially created, apply_testing: Application for adaptation, testing: Adaptation testing in progress, packaged: Tested and packaged,\r\nVerifying : Under review, Launched: Online, Removed: Removed
	Note            string    `gorm:"column:note;NOT NULL"`                                // Remarks
	Online          int       `gorm:"column:online;default:0;NOT NULL"`                    // Officially launched, 1: Officially launched, 0: Not officially launched, default is 1, when the game versions are incompatible, this field needs special configuration
	CreateAt        time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL"` // Creation time
	UpdateAt        time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL"` // Last modified time
}

func (m *AppVersion) TableName() string {
	return "app_version"
}
