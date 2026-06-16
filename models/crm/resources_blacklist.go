package crmModel

import "time"

type BlacklistResource struct {
	ID        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"` // 主键ID
	RID       int64     `gorm:"column:rid;type:varchar(64);NOT NULL"` // 资源编码
	TID       int64     `gorm:"column:tid;default:0;NOT NULL"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

// TableName 指定表名
func (b *BlacklistResource) TableName() string {
	return "resource_blacklist"
}
