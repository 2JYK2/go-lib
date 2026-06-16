package commonModel

import (
	"time"
)

// 容器黑名单
type ConfResourceBlacklist struct {
	ID       int64         `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Rid      int64         `gorm:"column:rid;NOT NULL" json:"rid"`
	Deleted  int           `gorm:"column:deleted;default:0;NOT NULL" json:"deleted"`
	EffectAt time.Time     `gorm:"column:effect_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"effect_at"` // 加入时间
	DeleteAt time.Time     `gorm:"column:delete_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"delete_at"` // 删除时间
	CreateAt time.Time     `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"` // 创建时间
	UpdateAt time.Time     `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"` // 最后更新时间
	Channel  CreateChannel `gorm:"column:channel;NOT NULL" json:"channel"`                               // 来源
}

func (m *ConfResourceBlacklist) TableName() string {
	return "conf_resource_blacklist"
}

type CreateChannel string

const (
	ManagerPortalCreate CreateChannel = "ManagerPortalCreate"
	CheckFailCreate     CreateChannel = "CheckFailCreate"
)
