package libModel

type OsType string

const (
	OsTypeWindows OsType = "windows"
	OsTypeAndroid OsType = "android"
	OsTypeLinux   OsType = "linux"
)

const (
	ContainerListeningPortPrefix = 10100
	ContainerIDPrefix            = 100
)

type ResourceType string

const (
	HOST      ResourceType = "host"
	CONTAINER ResourceType = "container"
	BENCHMARK ResourceType = "benchmark"
)

const (
	DataNoDelete = 0
	DataIsDelete = 1
)
const (
	CreateGidLimit = 10000
	FindGidBatch   = 1000
)
