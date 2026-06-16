package crmModel

type ResourceSshKeyRel struct {
	Id       int64 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Rid      int64 `gorm:"column:rid;default:0;NOT NULL" json:"rid"`
	SshKeyId int64 `gorm:"column:ssh_key_id;default:0;NOT NULL" json:"ssh_key_id"`
	Tid      int64 `gorm:"column:tid;default:0;NOT NULL" json:"tid"`
}

func (m *ResourceSshKeyRel) TableName() string {
	return "resource_ssh_key_rel"
}
