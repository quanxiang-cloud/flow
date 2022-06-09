package logger

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

func new() (*gorm.DB, error) {

	newLogger := NewGorm(Logger, lg.Config{
		SlowThreshold: time.Second,
		LogLevel:      lg.Info,
		Colorful:      true,
	})

	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
				"root", "uyWxtvt6gCOy3VPLB3rTpa0rQ", "127.0.0.1:3306", "config_center"),
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

	return db, nil
}

func TestGormLogger(t *testing.T) {
	err := New(&Config{
		Level:           -1,
		Development:     false,
		Sampling:        Sampling{100, 100},
		OutputPath:      []string{"stderr"},
		ErrorOutputPath: []string{"stderr"},
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	db, err := new()
	if err != nil {
		t.Fatal(err)
		return
	}
	var data interface{}
	err = db.Table("app").Find(&data).
		Error
	if err != nil {
		t.Fatal(err)
		return
	}
}
