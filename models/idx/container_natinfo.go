package idxModel

import (
	"time"
)

// 容器网络映射表
type ContainerNatInfo struct {
	Cid      int64     `gorm:"column:cid;type:bigint(20);primary_key;default:0" json:"cid"`
	NatIp    string    `gorm:"column:nat_ip;type:varchar(20);comment:nat_ip_v4;NOT NULL" json:"nat_ip"`
	NatPort  int       `gorm:"column:nat_port;type:int(11);default:0;comment:外网端口;NOT NULL" json:"nat_port"`
	CreateAt time.Time `gorm:"column:create_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"`
	UpdateAt time.Time `gorm:"column:update_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"`
}

func (m *ContainerNatInfo) TableName() string {
	return "container_natinfo"
}
