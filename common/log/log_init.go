package log

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* Usage:
 * Default log file:
 * log.Run.Debug("I am a Message", zap.String("key", value), zap.Int64("intKey", intValue))
 * Log by name:
 * log.GetLogger("loggerName").Debug("I am a Message", zap.String(key, value), zap.Int64(intKey, intValue))
 * Log by level:
 * log.Run.Log(zapcore.DebugLevel, "I am a Message", zap.String(key, value), zap.Int64(intKey, intValue))
 * Log by error:
 * log.Error.ErrorIf("I am a Message", error, zap.String(key, value), zap.Int64(intKey, intValue))
 */

const (
	// Basic log level logger name
	RunLogName    string = "running" // DefaultLogName must be configured in the config file, same as below
	ErrorLogName  string = "error"
	StatisLogName string = "statis"
)

var (
	runObj    *Logger
	errorObj  *Logger
	statisObj *Logger
	_loggers  map[string]*Logger
)

type Logger struct {
	*zap.Logger
}

func InitLog(logConfigs LogConfigs, modelName string) {
	_loggers = make(map[string]*Logger)
	for _, cf := range logConfigs.LogConfigs {
		initLog(cf, zapcore.Level(logConfigs.LowestLevel), zapcore.Level(logConfigs.StackLevel), modelName)
	}
	runObj = _loggers[RunLogName]
	statisObj = _loggers[StatisLogName]
	errorObj = _loggers[ErrorLogName]
}

func initLog(conf *LogConfig, lowestLevel zapcore.Level, stackLevel zapcore.Level, modelName string) {
	encoderConf := zapcore.EncoderConfig{
		TimeKey:  T,
		LevelKey: L,
		//FunctionKey:    Method,
		MessageKey:     Msg,
		StacktraceKey:  Stack,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	}

	// Create the encoder
	//bigDataEncoder := NewBigDataEncoder(encoderConf)
	//newJsonDataEncoder := NewJSONEncoder(encoderConf)

	// set lumberJack
	lumberJack := newLumberjackLogger(conf, modelName)

	var cores []zapcore.Core
	cores = append(cores,
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf), zapcore.AddSync(lumberJack),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= lowestLevel })))

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the cores together.
	core := zapcore.NewTee(cores...)
	// From a zapcore.Core, it's easy to construct a Logger.
	_loggers[conf.Filename] = &Logger{
		zap.New(core,
			zap.Fields(zap.String(ModelName, modelName)),
			zap.AddCaller(),
			zap.AddStacktrace(zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= stackLevel })),
		),
	}
	return
}

func newLumberjackLogger(conf *LogConfig, modelName string) *lumberjack.Logger {
	fileName := conf.LogPath + "/" + modelName + "_" + conf.Filename + ".log"
	logger := &lumberjack.Logger{
		Filename:   conf.LogPath + "/" + modelName + "_" + conf.Filename + ".log",
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups, // backup
		MaxAge:     conf.MaxAge,     // days
		Compress:   conf.Compress,   // disabled by default
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			_, err := os.Stat(fileName)
			if os.IsNotExist(err) {
				// create file
				if err = logger.Rotate(); err != nil {
					fmt.Println("Create File Err", err, fileName)
				}
			}
		}
	}()
	return logger
}
