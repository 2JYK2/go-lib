package log

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestExportM(t *testing.T) {
	logConfigs := LogConfigs{
		LowestLevel: 0,
		StackLevel:  0,
		LogConfigs: []*LogConfig{
			{
				Filename:   "running",
				LogPath:    "./log",
				MaxSize:    5,
				MaxBackups: 3,
				MaxAge:     30,
				Compress:   true,
			},
			{
				Filename:   "error",
				LogPath:    "./log",
				MaxSize:    5,
				MaxBackups: 3,
				MaxAge:     30,
				Compress:   true,
			},
		},
	}

	InitLog(logConfigs, "IDX")

	test()
}

func test() {
	ctx := context.Background()
	EmailErrSave(ctx, "测试",
		123,
		"123",
		fmt.Errorf("错误"),
		zap.Int(SID, 123),
		zap.Int(AppID, 1111),
		zap.Any("jyk", "{}"))
	GetLoggerByName(nil, RunLogName).Warn("test", zap.Error(fmt.Errorf("xerr")))

	time.Sleep(1 * time.Minute)
}
