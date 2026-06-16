package crmModel

import "time"

type ResourceSshKey struct {
	Id       int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Tid      int64     `gorm:"column:tid;default:0;NOT NULL" json:"tid"`
	Name     string    `gorm:"column:name;NOT NULL" json:"name"`
	SshKey   string    `gorm:"column:ssh_key;NOT NULL" json:"ssh_key"`
	CreateBy string    `gorm:"column:create_by;default:0;NOT NULL" json:"create_by"`
	CreateAt time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"`
	UpdateBy string    `gorm:"column:update_by;default:0;NOT NULL" json:"update_by"`
	UpdateAt time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"`
}

func (m *ResourceSshKey) TableName() string {
	return "resource_ssh_key"
}

const (
	CmdAddSshKey = `echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDo0EmFOb4rePwB/NA5j3d/x4Vq5ETfYHRy9CHiOv+QNWxNKDicB0CaklH/CD/OxuatVi449EqLl9FzJRW13NL4hjQtAK6lA3C7r9FlHoAwUUyGgWTdKov+9WoaZ11xq0O9rOXlsao++Jc2Pji1DlEs2clhWOfpkvsmgRW/TjoxqUw/greyFYBytsTuUX/Le5WFczEFvW+KKW6mKfK/h8HjfhsOA/zUyw/uHiOj49JTFEmhmSpE0ElQl1HH3ovk2Zhm2aYXqpZEmp/HBlmz1t/UrGWGo/rB7Aw46H2OK9lO8aoB+oWW3uIhLez+Y25WpyNYGl1kui2Qct4p2NifmxsgKycWUH1MdxEB2x6WmbILtI6ujIeu77yNejv0qf3E82BixEQvIjKLly6OeNytkxdS5j/SfhKURGmjmsO1mIwfmUXH20WXHv49qgAhoMLxocdm5l0O3LfGXNfJUvgk8+s1593Z8ESzO5QpQf2xwRCOdCJMFf/Azuvi/jotX1HPknM= yx@yx-System-Product-Name" >> ~/.ssh/authorized_keys`
	CmdDelSshKey = `sed -i '\|ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDo0EmFOb4rePwB/NA5j3d/x4Vq5ETfYHRy9CHiOv+QNWxNKDicB0CaklH/CD/OxuatVi449EqLl9FzJRW13NL4hjQtAK6lA3C7r9FlHoAwUUyGgWTdKov+9WoaZ11xq0O9rOXlsao++Jc2Pji1DlEs2clhWOfpkvsmgRW/TjoxqUw/greyFYBytsTuUX/Le5WFczEFvW+KKW6mKfK/h8HjfhsOA/zUyw/uHiOj49JTFEmhmSpE0ElQl1HH3ovk2Zhm2aYXqpZEmp/HBlmz1t/UrGWGo/rB7Aw46H2OK9lO8aoB+oWW3uIhLez+Y25WpyNYGl1kui2Qct4p2NifmxsgKycWUH1MdxEB2x6WmbILtI6ujIeu77yNejv0qf3E82BixEQvIjKLly6OeNytkxdS5j/SfhKURGmjmsO1mIwfmUXH20WXHv49qgAhoMLxocdm5l0O3LfGXNfJUvgk8+s1593Z8ESzO5QpQf2xwRCOdCJMFf/Azuvi/jotX1HPknM= yx@yx-System-Product-Name|d' ~/.ssh/authorized_keys`
)
