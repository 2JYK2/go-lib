package azlog

type Config struct {
	LogPath    string
	Filename   string
	Level      int8
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	LocalTime  bool
}
