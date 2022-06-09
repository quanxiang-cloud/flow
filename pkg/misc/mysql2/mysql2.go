package mysql2

import (
	"fmt"
	"time"

	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

// Config 配置
type Config struct {
	Host     string
	DB       string
	User     string
	Password string
	Log      bool

	MaxIdleConns int
	MaxOpenConns int

	dsn string
}

// DSN code
const (
	// DSNDEFAULT utf-8
	DSNDEFAULT = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"

	// DSNUTF8MB4 utf-8 mb4
	DSNUTF8MB4 = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

// String String
func (c *Config) String() string {
	if c.dsn == "" {
		c.dsn = DSNDEFAULT
	}
	return c.string(c.dsn)
}

// SetDSN SetDSN
func (c *Config) SetDSN(dsn string) {
	c.dsn = dsn
}

func (c *Config) string(format string) string {
	return fmt.Sprintf(format, c.User, c.Password, c.Host, c.DB)
}

// New 创建数据库连接
func New(config Config, log *zap.SugaredLogger) (*gorm.DB, error) {
	var newLogger lg.Interface
	if config.Log {
		newLogger = logger.NewGorm(log, lg.Config{
			SlowThreshold: time.Second,
			LogLevel:      lg.Info,
			Colorful:      true,
		})
	}
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       config.String(),
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}),
		&gorm.Config{
			Logger: newLogger,
		})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 20
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	return db, nil
}

// DeleteCallBack DeleteCallBack
func DeleteCallBack(db *gorm.DB) {
	db.Callback().Query().Before("*").Register("plugin:before_delete", func(db *gorm.DB) {
		if db.Statement.Schema == nil {
			return
		}
		if _, ok := db.Statement.Schema.FieldsByDBName["deleted_at"]; ok {
			db = db.Where("deleted_at = 0")
		}
	})
}
