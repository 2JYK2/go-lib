package log

type LogConfigs struct {
	LowestLevel int8         `yaml:"LowestLevel"` // level(Debug:-1, Info:0, Warn:1, Error:2, Dpanic:3, Panic:4, Fatal:5)
	StackLevel  int8         `yaml:"StackLevel"`  // The lowest log level to print stack information
	LogConfigs  []*LogConfig `yaml:"NamedLoggers"`
}

// Log related configuration
type LogConfig struct {
	// file paths
	LogPath string `json:"LogPath" yaml:"LogPath"`
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `json:"Filename" yaml:"Filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"MaxSize" yaml:"MaxSize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"MaxAge" yaml:"MaxAge"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"MaxBackups" yaml:"MaxBackups"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"Compress" yaml:"Compress"`
}

const (
	T             string = "time"
	L             string = "level"
	Line          string = "Line"
	Method        string = "Method"
	Msg           string = "Msg"
	Stack         string = "Stack"
	ModelName     string = "MN"
	TraceID       string = "TraceID"
	SID           string = "SID"
	SIDS          string = "SIDS"
	TID           string = "TID"
	AppID         string = "AppID"
	AppVersion    string = "AppVersion"
	UID           string = "UID"
	CID           string = "CID"
	CIDS          string = "CIDS"
	Action        string = "Action"
	Api           string = "Api"
	Success       string = "Success"
	Code          string = "Code"
	Err           string = "Err"
	EvtAt         string = "EvtAt"
	CostMs        string = "CostMs"
	CostUs        string = "CostUs"
	ReqUrl        string = "ReqUrl"
	ReqData       string = "ReqData"
	RespData      string = "RespData"
	Alarm         string = "Alarm"
	Ti            string = "Ti"
	RedisKey      string = "RedisKey"
	RedisValue    string = "RedisValue"
	RegionCode    string = "RegionCode"
	RegionCodeS   string = "RegionCodeS"
	ClientModel   string = "ClientModel"
	ClientType    string = "ClientType"
	Fingerprint   string = "Fingerprint"
	FingerprintS  string = "FingerprintS"
	Spec          string = "Spec"
	WholeSaleTid  string = "WholeSaleTid"
	SessionStatus string = "SessionStatus"
	ReqIp         string = "ReqIp"
	ReqTime       string = "ReqTime"
	UserAgent     string = "UserAgent"
	Header        string = "Header"
	Referer       string = "Referer"
	RespCode      string = "RespCode"
	StartTime     string = "StartTime"
	EndTime       string = "EndTime"
	Solution      string = "Solution"
	EmailAlarm    int    = 1
	ShortMsgAlarm int    = 2
	PhoneAlarm    int    = 3
)

const (
	IDX          string = "IDX"
	REGISTER     string = "REGISTER"
	AppAdapter   string = "AppAdapter"
	CRM          string = "CRM"
	Pkg          string = "Pkg"
	PsJob        string = "PsJob"
	GpuScheduler string = "GpuScheduler"
)
