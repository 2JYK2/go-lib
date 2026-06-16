package log

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

func NewGormLogger(config logger.Config) *gormLogger {
	l := &gormLogger{Config: config}
	l.infoStr = "%s [info] "
	l.warnStr = "%s [warn] "
	l.errStr = "%s [error] "
	l.traceStr = "%s [%.3fms] [rows:%v] %s"
	l.traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	l.traceErrStr = "%s %s [%.3fms] [rows:%v] %s"
	return l
}

type gormLogger struct {
	logger.Config
	ctx                                 context.Context
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func (g *gormLogger) Printf(s string, i ...interface{}) {
	l := GetLogger(RunLogName, g.ctx, 6)
	if l != nil {
		l.Info(fmt.Sprintf(strings.ReplaceAll(s, "\n", " "), i...))
	}
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *g
	newlogger.LogLevel = level
	return &newlogger
}

func (g *gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel >= logger.Info {
		g.ctx = ctx
		g.Printf(g.infoStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

func (g *gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel >= logger.Warn {
		g.ctx = ctx
		g.Printf(g.warnStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

func (g *gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if g.LogLevel >= logger.Error {
		g.ctx = ctx
		g.Printf(g.errStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
	}
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}
	g.ctx = ctx
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !g.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			g.Printf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Printf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		if rows == -1 {
			g.Printf(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Printf(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case g.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			g.Printf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			g.Printf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
