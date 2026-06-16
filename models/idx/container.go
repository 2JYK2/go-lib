package idxModel

import (
	"time"

	libModel "github.com/2JYK2/go-lib/models"
)

// container info
type Container struct {
	ID                   int64           `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	RegionCode           string          `gorm:"column:region_code;NOT NULL"`                    // Region code
	Fingerprint          string          `gorm:"column:fingerprint;NOT NULL"`                    // fingerprint
	TID                  int64           `gorm:"column:tid;default:0;NOT NULL"`                  // Belonging tenant ID, i.e. mining ID
	SID                  int64           `gorm:"column:sid;default:0;NOT NULL"`                  // Session ID
	OsType               libModel.OsType `gorm:"column:os_type;default:windows;NOT NULL"`        // Container system type, windows, android
	LocalIP              string          `gorm:"column:local_ip;NOT NULL"`                       // 本机IP
	LocalPort            int             `gorm:"column:local_port;default:0;NOT NULL"`           // 本地音视频流监听端口
	ContainerType        ContainerType   `gorm:"column:container_type;default:0;NOT NULL"`       // Container Type
	ArsMgrStartTs        int64           `gorm:"column:ars_mgr_start_ts;default:0;NOT NULL"`     // arsdog latest restart timestamp
	ContainerStartTs     int64           `gorm:"column:container_start_ts;default:0;NOT NULL"`   // Container latest start timestamp
	ContainerHealthyTs   int64           `gorm:"column:container_healthy_ts;default:0;NOT NULL"` // Heartbeat latest start timestamp
	ContainerConnIdxAddr string          `gorm:"column:container_conn_idx_addr;NOT NULL"`        // ars module long connection indexer node address IP:port
	Status               ContainerStatus `gorm:"column:status;NOT NULL"`                         // Current status
	StatusChangeAt       int64           `gorm:"column:status_change_at;NOT NULL"`               // Current status modification timestamp
	Online               bool            `gorm:"column:online;default:1;NOT NULL"`               // 0 未上架 1 已上架
	Healthy              bool            `gorm:"column:healthy;default:1;NOT NULL"`              // Healthy status, 1: available, 0: unavailable, heartbeat lost
	Lock                 int             `gorm:"column:lock;default:0;NOT NULL"`                 // Whether in maintenance state, 1: locked, 0: unlocked
	BizType              int             `gorm:"column:biz_type;default:0;NOT NULL"`             // 业务类型（0云游戏，1云桌面2.云游戏&&云桌面)
	WholeSaleTid         int64           `gorm:"column:whole_sale_tid;default:0;NOT NULL"`       // 订购容器的租户id
	Pl                   string          `gorm:"column:pl;NOT NULL"`                             // 渲染等级
	Location             string          `gorm:"column:location;NOT NULL"`                       // 城市
	UpdateChecker        string          `gorm:"column:update_checker;NOT NULL"`                 // Data update operation optimistic lock
	ErrorMsg             string          `gorm:"column:error_msg;default:0;NOT NULL"`            // Start error message
	CreateAt             time.Time       `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	UpdateAt             time.Time       `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL"`
	SortWeight           int64           `gorm:"-"`
}

func (m *Container) TableName() string {
	return "container"
}

type RoomType int32

type ContainerStatus string

const (
	ContainerStatusError        ContainerStatus = "Error"
	ContainerStatusInitializing ContainerStatus = "Initializing"
	ContainerStatusIdle         ContainerStatus = "Idle"
	ContainerStatusPreparing    ContainerStatus = "Preparing"
	ContainerStatusGaming       ContainerStatus = "Gaming"
	ContainerStatusToBeReleased ContainerStatus = "ToBeReleased"
	ContainerStatusReleasing    ContainerStatus = "Releasing"
)

type CheckStatus string

const (
	WaitSpec  CheckStatus = "wait_spec"
	CheckWait CheckStatus = "check_wait"
	CheckPass CheckStatus = "check_pass"
)

type ContainerType int

const (
	ContainerAgent ContainerType = 0
	Iaas           ContainerType = 1
)

const (
	ArmLocalPort = 6001
	X86LocalPort = 22001
)
