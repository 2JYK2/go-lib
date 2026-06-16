package idxModel

import (
	"time"
)

// Session service resource
type Session struct {
	AutoID               int64         `gorm:"column:auto_id;primary_key;AUTO_INCREMENT" json:"AutoID"`
	SID                  int64         `gorm:"column:sid" json:"SID"`
	SessionType          int           `gorm:"column:session_type;NOT NULL" json:"sessionType"`
	Spec                 string        `gorm:"column:spec;NOT NULL" json:"specName"`                               // container type
	RegionCode           string        `gorm:"column:region_code;NOT NULL" json:"regionCode"`                      // region code
	ContainerID          int64         `gorm:"column:container_id;default:0;NOT NULL" json:"containerID"`          // container ID
	Status               SessionStatus `gorm:"column:status;NOT NULL" json:"status"`                               // current status
	TID                  int64         `gorm:"column:tid;NOT NULL" json:"providerID"`                              // provider id
	UID                  string        `gorm:"column:uid;NOT NULL" json:"UID"`                                     // user id
	AppID                int64         `gorm:"column:app_id;default:0;NOT NULL" json:"appID"`                      // game ID running on the container, 0 if no game
	AppVersion           int64         `gorm:"column:app_version;default:0;NOT NULL" json:"appVersion"`            // game version
	ClientModel          string        `gorm:"column:client_model;NOT NULL" json:"clientModel"`                    // terminal device model identifier
	ClientType           int           `gorm:"column:client_type;default:0;NOT NULL" json:"clientType"`            // current service client type
	ClientIP             string        `gorm:"column:client_ip;NOT NULL" json:"clientIP"`                          // client IP
	ClientConnectionTime int64         `gorm:"column:client_connection_time;NOT NULL" json:"clientConnectionTime"` // client connection time
	NetType              NetType       `gorm:"column:net_type;NOT NULL" json:"netType"`                            // client network identifier
	TermDevID            string        `gorm:"column:term_dev_id;NOT NULL" json:"termDevID"`                       // terminal device unique identifier
	ByTest               bool          `gorm:"column:by_test;default:0;NOT NULL" json:"byTest"`                    // is test scenario
	ErrorCode            int64         `gorm:"column:error_code;default:0;NOT NULL" json:"errorCode"`              // error code
	ErrorMsg             string        `gorm:"column:error_msg;NOT NULL" json:"errorMsg"`                          // error message
	QuitReason           string        `gorm:"column:quit_reason;NOT NULL" json:"quitReason"`                      // event code
	CallbackEvtKey       string        `gorm:"column:callback_evt_key;NOT NULL" json:"callbackEvtKey"`             // event callback private parameter
	NoInputTimeouts      int           `gorm:"column:no_input_timeouts;NOT NULL" json:"noInputTimeouts"`           // no input timeout
	StartAppParam        string        `gorm:"column:start_app_param;NOT NULL" json:"startAppParam"`               // start game passthrough parameter
	Language             string        `gorm:"column:language;NOT NULL" json:"language"`                           // server language
	Timezone             string        `gorm:"column:timezone;NOT NULL" json:"timezone"`                           // server timezone
	BizType              int           `gorm:"column:biz_type;NOT NULL" json:"bizType"`                            // business type
	FirstFrameTime       int64         `gorm:"column:first_frame_time;NOT NULL" json:"firstFrameTime"`             // first frame time
	IsWholeSale          string        `gorm:"column:is_whole_sale;NOT NULL"`                                      // is long-term rental mode
	ArchiveMode          int           `gorm:"column:archive_mode;default:0;NOT NULL"`                             // archive handling mode
	ArchiveID            int64         `gorm:"column:archive_id;default:0;NOT NULL"`                               // archive ID
	AbnormalUpload       bool          `gorm:"column:abnormal_upload;default:0;NOT NULL"`                          // restore archive in abnormal scenarios

	Checker              string `gorm:"column:checker;NOT NULL" json:"checker"`                              // modification verification
	ClientScreen         string `gorm:"column:client_screen;NOT NULL" json:"clientScreen"`                   // client screen
	AppGraphicsLevel     string `gorm:"column:app_graphics_level;NOT NULL" json:"appGraphicsLevel"`          // game launch graphics level identifier
	AppSupportDevices    string `gorm:"column:app_support_devices;NOT NULL" json:"appSupportDevices"`        // game support devices list
	AppVideoQualityLevel string `gorm:"column:app_video_quality_level;NOT NULL" json:"appVideoQualityLevel"` // current session video stream quality level
	PlayUrl              string `gorm:"column:play_url;NOT NULL" json:"playUrl"`                             // RTC access address
	PlayToken            string `gorm:"column:play_token;NOT NULL" json:"playToken"`                         // RTC access token

	UpdateAt time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP" json:"updateAt"`
	CreateAt time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP" json:"createAt"`
}

// TableName get sql table name.
func (m *Session) TableName() string {
	return "session"
}

func (s *Session) InWorkingStatus() bool {
	return s.Status == SessionStatus_InService || s.Status == SessionStatus_Prepared
}

func (s *Session) InFinishedStatus() bool {
	return s.Status == SessionStatus_Finished || s.Status == SessionStatus_Error
}

type SessionStatus string

const (
	SessionStatus_Error     SessionStatus = "error"
	SessionStatus_Preparing SessionStatus = "preparing"
	SessionStatus_Prepared  SessionStatus = "prepared"
	SessionStatus_InService SessionStatus = "in-service"
	SessionStatus_Finished  SessionStatus = "finished"
	SessionStatus_OnQueue   SessionStatus = "on-queue"
	SessionStatus_NoIdle    SessionStatus = "no-idle"
	SessionStatus_HangUp    SessionStatus = "hangup"
)

type NetType string

const (
	Local   NetType = "local"
	Default NetType = "default"
)
