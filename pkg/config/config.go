package config

import (
	"github.com/quanxiang-cloud/flow/pkg/misc/client"
	"github.com/quanxiang-cloud/flow/pkg/misc/kafka"
	"github.com/quanxiang-cloud/flow/pkg/misc/mysql2"
	"github.com/quanxiang-cloud/flow/pkg/misc/redis2"
	"gopkg.in/yaml.v2"
	"os"
	"time"

	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
)

// Config 全局配置对象
var Config *Configs

// Configs 总配置结构体
type Configs struct {
	Model       string        `yaml:"model"`
	Port        string        `yaml:"port"`
	Mysql       mysql2.Config `yaml:"mysql"`
	Log         logger.Config `yaml:"log"`
	InternalNet client.Config `yaml:"internalNet"`
	Redis       redis2.Config `yaml:"redis"`
	Kafka       kafka.Config  `yaml:"kafka"`
	APIHost     APIHost       `yaml:"api"`
}

// HTTPServer http服务配置
type HTTPServer struct {
	Port              string        `yaml:"port"`
	ReadHeaderTimeOut time.Duration `yaml:"readHeaderTimeOut"`
	WriteTimeOut      time.Duration `yaml:"writeTimeOut"`
	MaxHeaderBytes    int           `yaml:"maxHeaderBytes"`
}

// APIHost api host
type APIHost struct {
	OrgHost           string `yaml:"orgHost" validate:"required"`
	GoalieHost        string `yaml:"goalieHost" validate:"required"`
	FormHost          string `yaml:"formHost" validate:"required"`
	AppCenterHost     string `yaml:"appCenterHost" validate:"required"`
	MessageCenterHost string `yaml:"messageCenterHost" validate:"required"`
	StructorHost      string `yaml:"structorHost" validate:"required"`
	DispatcherHost    string `yaml:"dispatcherHost" validate:"required"`
	ProcessHost       string `yaml:"processHost" validate:"required"`
	PolyAPIHost       string `yaml:"polyAPIHost" validate:"required"`
}

// Init 初始化
func Init(configPath string) error {
	if configPath == "" {
		configPath = "../configs/configs.yml"
	}
	Config = new(Configs)
	err := read(configPath, Config)
	if err != nil {
		return err
	}
	return nil
}

// read 读取配置文件
func read(yamlPath string, v interface{}) error {
	// Read config file
	buf, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	return nil
}
