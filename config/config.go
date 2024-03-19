package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Server         ServerConfig    `mapstructure:"api"`
	Logger         LoggerConfig    `mapstructure:"logger"`
	Authentication *Authentication `mapstructure:"authentication"`
	Database       Database        `mapstructure:"database"`
	Mysql          *Mysql          `mapstructure:"mysql"`
	Sqlite         *Sqlite         `mapstructure:"sqlite"`
	Env            string          `mapstructure:"env"`
}

type ServerConfig struct {
	Port  string `mapstructure:"port"  default:"9090"`
	Debug bool   `mapstructure:"debug" default:"false"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type Authentication struct {
	ApiKeys     []string `mapstructure:"api_keys"`
	Users       []string `mapstructure:"users"`
	SecretKey   string   `mapstructure:"secret_key"`
	ExpiredTime int64    `mapstructure:"expired_time"`
	CookieName  string   `mapstructure:"cookie_name"`
}

type Mysql struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"database_name"`
}

func (m Mysql) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DatabaseName,
	)
}

type Sqlite struct {
	Path string `mapstructure:"path"`
}

func (s Sqlite) GetDSN() string {
	return s.Path
}

type Database struct {
	Debug             bool   `mapstructure:"debug"`
	DBType            string `mapstructure:"db_type"`
	MaxLifetime       int    `mapstructure:"max_lifetime"`
	MaxOpenConns      int    `mapstructure:"max_open_conns"`
	MaxIdleConns      int    `mapstructure:"max_idle_conns"`
	TablePrefix       string `mapstructure:"table_prefix"`
	EnableAutoMigrate bool   `mapstructure:"enable_auto_migrate"`
}

var (
	_, b, _, _        = runtime.Caller(0)
	basePath          = filepath.Dir(b) //get absolute directory of current file
	defaultConfigFile = basePath + "/local.yaml"
	v                 = viper.New()
	appConfig         AppConfig
)

func Load() {
	var configFile string
	if configFile = os.Getenv("CONFIG_PATH"); len(configFile) == 0 {
		configFile = defaultConfigFile
	}

	if err := loadConfigFile(configFile); err != nil {
		panic(err)
	}

	if err := scanConfigFile(&appConfig); err != nil {
		panic(err)
	}
}

func loadConfigFile(configFile string) error {
	configFileName := filepath.Base(configFile)
	configFilePath := filepath.Dir(configFile)

	v.AddConfigPath(configFilePath)
	v.SetConfigName(strings.TrimSuffix(configFileName, filepath.Ext(configFileName)))
	v.AutomaticEnv()

	return v.ReadInConfig()
}

func scanConfigFile(config any) error {
	return v.Unmarshal(&config)
}

func GetAppConfig() *AppConfig {
	return &appConfig
}
