package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/2JYK2/go-lib/common/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// // GCWBossDB library operation object
// var GCWBossDB *gorm.DB
//
// // CenterControllerDB center-controller operation object
// var CenterControllerDB *gorm.DB
// Custom logger for gorm
type GormLogger struct {
	SlowThreshold time.Duration
}

func newGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		SlowThreshold: slowThreshold,
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return l
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {

}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {

}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	// Gets the SQL statement and the number of returns
	sql, rows := fc()

	log.DebugSave(ctx, fmt.Sprintf("SQL INFO, | sql=%v, rows=%v, elapsed=%v", sql, rows, elapsed))

	// Gorm error printed
	if err != nil && err != gorm.ErrRecordNotFound && !strings.Contains(err.Error(), "context canceled") {
		log.EmailErrSave(ctx, fmt.Sprintf("SQL ERROR, | sql=%v, rows=%v, elapsed=%v", sql, rows, elapsed), 500, "", err)
	}
	// Slow query log
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.WarnSave(ctx, fmt.Sprintf("Database Slow Log, | sql=%v, rows=%v, elapsedTs=%v", sql, rows, elapsed), err)
	}
}

func InitMysqlDao(mysqlDbConfig *MysqlDBConfig) *gorm.DB {
	dsnRW := mysqlDbConfig.WRITE.DBUrl(mysqlDbConfig.ConnParam)
	dsnR := mysqlDbConfig.READ.DBUrl(mysqlDbConfig.ConnParam)

	Sql, err := gorm.Open(mysql.Open(dsnRW), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         newGormLogger(time.Second * time.Duration(mysqlDbConfig.SlowThresholdUs)),
	})

	if Sql == nil {
		return nil
	}

	err = Sql.Use(dbresolver.Register(dbresolver.Config{
		// use `db2` as sources, `db3`, `db4` as replicas
		Sources:           []gorm.Dialector{mysql.Open(dsnRW)},
		Replicas:          []gorm.Dialector{mysql.Open(dsnR)},
		TraceResolverMode: true},
	).SetConnMaxIdleTime(time.Duration(mysqlDbConfig.ConnMaxIdleTime) * time.Second).
		SetMaxIdleConns(mysqlDbConfig.MaxIdleConns).
		SetMaxOpenConns(mysqlDbConfig.MaxOpenConns).
		SetConnMaxLifetime(time.Duration(mysqlDbConfig.ConnMaxLifetime) * time.Second))
	if err != nil {
		panic(fmt.Sprintf("open crm db fail  err=%v", err))
	}

	return Sql
}
