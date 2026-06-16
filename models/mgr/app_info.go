package mgrModel

import (
	libModel "github.com/2JYK2/go-lib/models"
	"time"
)

// Application Information
type AppInfo struct {
	AppID                int64           `gorm:"column:app_id;primary_key;AUTO_INCREMENT" json:"app_id"` // Unique identifier of the game, globally unique, with an ID less than 1000, indicating that it is a rendering engine working software
	AppName              string          `gorm:"column:app_name;NOT NULL" json:"app_name"`               // apply name
	Tid                  int64           `gorm:"column:tid;default:0;NOT NULL" json:"tid"`               // Tenant, if blank, represents multi tenant sharing
	Spec                 string          `gorm:"column:spec;NOT NULL" json:"spec"`                       //
	PkgName              string          `gorm:"column:pkg_name" json:"pkg_name"`
	LauncherName         string          `gorm:"column:launcher_name;NOT NULL" json:"launcher_name"`       // The launcher name of the app required for Android game launch
	PlatformType         libModel.OsType `gorm:"column:platform_type;NOT NULL" json:"platform_type"`       // Application platform type, including: Windows, Android
	PlatformSubtype      string          `gorm:"column:platform_subtype;NOT NULL" json:"platform_subtype"` // Platform subtype
	AppType              string          `gorm:"column:app_type;NOT NULL" json:"app_type"`                 // Game type, GameOnline SteamLibrary
	SteamAppID           string          `gorm:"column:steam_app_id;NOT NULL" json:"steam_app_id"`         // Application Id of Steam Game
	InstallDir           string          `gorm:"column:install_dir;NOT NULL" json:"install_dir"`           // Relative directory for game deployment
	ConfigDir            string          `gorm:"column:config_dir;NOT NULL" json:"config_dir"`             // Game Configuration Directory
	ConfigMd5            string          `gorm:"column:config_md5;NOT NULL" json:"config_md5"`             // The game configuration file MD5 matches the configuration version
	SupportInputDevices  string          `gorm:"column:support_input_devices;NOT NULL" json:"support_input_devices"`
	RunMode              int             `gorm:"column:run_mode;default:0;NOT NULL" json:"run_mode"`                     // The virtualization instance modes supported by the game: 0: desktop single instance, 1: sandbox single instance, 2: sandbox multiple instances, 3: non sandbox process multiple instances
	WindowMode           int             `gorm:"column:window_mode;default:0;NOT NULL" json:"window_mode"`               // Game window mode, 0: fixed size does not support window resizing in operation, 1: window resizing in operation is supported
	IsFullScreen         int             `gorm:"column:is_full_screen;default:0;NOT NULL" json:"is_full_screen"`         // Does it support full screen display
	IsSupportArchive     int             `gorm:"column:is_support_archive;default:0;NOT NULL" json:"is_support_archive"` // Does it support archiving, default is 0, not supported
	Orientation          int             `gorm:"column:orientation;default:0;NOT NULL" json:"orientation"`               // Horizontal and vertical screens, 0 horizontal and 1 vertical screen
	EncConfig            string          `gorm:"column:enc_config" json:"enc_config"`
	DefaultGraphicsLevel string          `gorm:"column:default_graphics_level;NOT NULL" json:"default_graphics_level"` // Game default image quality level
	StartupTime          int             `gorm:"column:startup_time;default:10;NOT NULL" json:"startup_time"`          // Game startup duration
	CodingStrategy       int             `gorm:"column:coding_strategy;default:0;NOT NULL" json:"coding_strategy"`     // Default game encoding strategy [0: forbidden 1: fixed frame rate 2: fixed resolution 3: balanced mode]
	Deleted              int             `gorm:"column:deleted;default:0;NOT NULL" json:"deleted"`                     // Has it been deleted? 1: It has been deleted. After deleting the game, all versions will naturally become invalid
	Description          string          `gorm:"column:description;NOT NULL" json:"description"`                       // Game Description
	CreateAt             time.Time       `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"` // Creation time
	UpdateAt             time.Time       `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"` // last-modified
	WithArchive          int             `gorm:"column:with_archive;default:0;NOT NULL" json:"with_archive"`           // Whether to support archiving, default is 0, not supported
	ExeArgs              string          `gorm:"column:exe_args;default:0;NOT NULL" json:"exe_args"`                   //
	ConfigContent        string          `gorm:"column:config_content;default:0;NOT NULL" json:"config_content"`       //
}

func (m *AppInfo) TableName() string {
	return "app_info"
}

type AppRunMode int

const (
	AppRunMode_Desktop       AppRunMode = 0
	AppRunMode_SandboxSingle AppRunMode = 1
	AppRunMode_SandboxMulti  AppRunMode = 2
	AppRunMode_ProcessMulti  AppRunMode = 2
)
