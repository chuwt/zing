package config

import (
	"fmt"
	"github.com/chuwt/zing/client/redis"
	"github.com/chuwt/zing/json"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var Config = &struct {
	StrategyPath string `yaml:"strategy_path"`

	Strategy Strategy `yaml:"strategy"`

	Mysql MysqlCfg     `yaml:"mysql"`
	Redis redis.Config `yaml:"redis"`

	LogLevel string `yaml:"log_level" default:"info"`

	DebugApiKey string `yaml:"debug_api_key"`
}{}

func init() {
	if err := configor.Load(Config, getConfigPath("")); err != nil {
		panic(err)
	}
	log := initLog(Config.LogLevel)
	configBytes, _ := json.Json.Marshal(Config)
	log.Info("读取配置文件", zap.String("config", string(configBytes)))
}

func initLog(logLevel string) *zap.Logger {
	logConf := zap.NewProductionConfig()
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	logConf.EncoderConfig = encoder
	logConf.Encoding = "console"

	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		panic(err)
	}
	logConf.Level = zap.NewAtomicLevelAt(level)

	log, err := logConf.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(log)
	return log
}

type Strategy struct {
	Path       string `yaml:"path"`
	PythonPath string `yaml:"python_path"`
}

type MysqlCfg struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`

	MaxConn      int `yaml:"max_conn"`
	MaxIdleConn  int `yaml:"max_idle_conn"`
	ConnLifetime int `yaml:"conn_lifetime"`
}

func (cfg *MysqlCfg) DSN() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
}

var osPathSeparator = string(os.PathSeparator)

func getConfigPath(appName string) string {
	var configPath string
	// get from evn
	if appName != "" {
		configPath = os.Getenv(appName)
	}
	// find conf/app.yml from outer dir
	if configPath == "" || appName == "" {
		nowDir, err := os.Getwd()
		if err != nil {
			panic("can't get root dir")
		}
		for {
			dirs, err := ioutil.ReadDir(nowDir)
			if err != nil {
				panic("can't read root dirs")
			}
			if target := getConfDir(dirs); target != "" {
				return path.Join(path.Join(nowDir, target), "app.yaml")
			}
			pathList := strings.Split(nowDir, osPathSeparator)
			nowDir = strings.Join(pathList[:len(pathList)-1], osPathSeparator)
		}
	}
	return configPath
}

func getConfDir(dirPath []os.FileInfo) string {
	for _, dir := range dirPath {
		if dir.IsDir() && dir.Name() == "conf" {
			return dir.Name()
		}
	}
	return ""
}
