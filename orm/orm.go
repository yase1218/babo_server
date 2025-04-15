package orm

import (
	"babo/utility/zlog"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// //////////////////////////////////////////////////////////////////////////////////////////
// mysql manager
type MysqlManager struct {
	db *gorm.DB
}

var mgr MysqlManager

func Connect(addr, user, pwd, db_name string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, addr, db_name)
	if mgr.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, //禁用默认事务
		PrepareStmt:            true, //启用预编译
		Logger:                 NewGormLog(),
		//Logger: pMgr.GetDBLogger(),
	}); err != nil {
		zlog.Error("Connect mysql failed", zap.Error(err))
		return
	}

	db, _ := mgr.db.DB()
	db.SetMaxOpenConns(240)
	db.SetMaxIdleConns(24)
	//_db.SetConnMaxIdleTime()
	db.SetConnMaxLifetime(time.Minute * 6)

	return
}

func DB() *gorm.DB {
	return mgr.db
}

// ////////////////////////////////////////////////////////////////////////////////////////
// redirect gorm logger
type GormLog struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
}

func NewGormLog() *GormLog {
	return &GormLog{
		SlowThreshold: time.Second,
		LogLevel:      logger.Warn,
	}
}
func (l *GormLog) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *GormLog) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		zlog.Infof(str, args...)
	}
}

func (l *GormLog) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		zlog.Warnf(str, args...)
	}
}

func (l *GormLog) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		zlog.Errorf(str, args...)
	}
}

func (l *GormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)

	sql, rows := fc()

	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.Any("time", elapsed),
		zap.Int64("rows", rows),
	}

	if err != nil {
		// 忽略未找到的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// l.logger.Warn("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, zap.Error(err))
			zlog.Error("SQL Error", logFields...)
		}
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		zlog.Warn("SQL Slow Log", logFields...)
	}

	// gorm.Debug()model 记录所有 SQL 请求
	if l.LogLevel == logger.Info {
		zlog.Info("SQL", logFields...)
	}

	// 记录所有SQL请求
	//Debug("SQL Query : {%s}", logStr)
}
