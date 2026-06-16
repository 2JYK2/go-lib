package gpuSchedulerModel

import (
	"time"
)

// gpu info
type Gpu struct {
	ID               int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	RegionCode       string    `gorm:"column:region_code;NOT NULL" json:"region_code"`                     // region_code
	Sid              int64     `gorm:"column:sid;default:0;NOT NULL" json:"sid"`                           // 会话id
	HostID           int64     `gorm:"column:host_id;default:0;NOT NULL" json:"host_id"`                   // 宿主机id
	OsType           string    `gorm:"column:os_type;default:windows;NOT NULL" json:"os_type"`             // Container system type, Windows, Android
	LocalIP          string    `gorm:"column:local_ip;NOT NULL" json:"local_ip"`                           // 本机IP
	LocalPort        int       `gorm:"column:local_port;default:0;NOT NULL" json:"local_port"`             // 本地音视频流监听端口
	ComfyuiStartTs   int64     `gorm:"column:comfyui_start_ts;default:0;NOT NULL" json:"comfyui_start_ts"` // 容器最近一次启动时间戳
	ComfdStartTs     int64     `gorm:"column:comfd_start_ts;default:0;NOT NULL" json:"comfd_start_ts"`     // 心跳时间
	ComfdConnIdxAddr string    `gorm:"column:comfd_conn_idx_addr;NOT NULL" json:"comfd_conn_idx_addr"`     // The node address IP: port of the long connection between the ars module and the indexer
	StatusChangeAt   int64     `gorm:"column:status_change_at;default:0;NOT NULL" json:"status_change_at"` // 状态修改时间
	Status           GpuStatus `gorm:"column:status;NOT NULL" json:"status"`                               // current state
	Healthy          int       `gorm:"column:healthy;default:0;NOT NULL" json:"healthy"`                   // Health status, 1: available, 0: unavailable, heartbeat loss
	LastHealthyTs    int       `gorm:"column:last_healthy_ts;default:0;NOT NULL" json:"last_healthy_ts"`   // 最后一次心跳时间
	Lock             int       `gorm:"column:lock;default:0;NOT NULL" json:"lock"`                         // 是否锁定 0未锁定 !0 锁定
	GpuType          string    `gorm:"column:gpu_type;NOT NULL" json:"gpu_type"`                           // gpu type
	GpuMem           int       `gorm:"column:gpu_mem;default:0;NOT NULL" json:"gpu_mem"`                   // gpu 显存
	ImageVersion     string    `gorm:"column:image_version;NOT NULL" json:"image_version"`                 // image version
	UpdateChecker    string    `gorm:"column:update_checker;NOT NULL" json:"update_checker"`               // Optimistic locks for data update operations
	NetIP            string    `gorm:"column:net_ip;NOT NULL" json:"net_ip"`                               // 外网入口ip
	NetPort          int       `gorm:"column:net_port;default:0;NOT NULL" json:"net_port"`                 // 外网入口port
	CreateAt         time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"`
	UpdateAt         time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"`
}

func (m *Gpu) TableName() string {
	return "gpu"
}

type GpuStatus string

const (
	Initializing GpuStatus = "Initializing"
	Idle         GpuStatus = "Idle"
	Running      GpuStatus = "Running"
	Releasing    GpuStatus = "Releasing"
	Error        GpuStatus = "Error"
	Finished     GpuStatus = "finished"
)
