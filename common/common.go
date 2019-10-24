package common

import (
	"flag"
	"time"

	"github.com/jmoiron/sqlx"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	gocommon "github.com/liuhengloveyou/go-common"
	passportcommon "github.com/liuhengloveyou/passport/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ClientConfigStruct struct {
	PID           string `yaml:"pid"`
	Name          string `yaml:"name"`
	TaskServeAddr string `yaml:"task_serve_addr"`
	TaskType      string `yaml:"task_type"`
	Flow          int    `yaml:"flow"`
	LogDir        string `yaml:"log_dir"`
	LogLevel      string `yaml:"log_level"`
}

type ServeConfigStruct struct {
	PID      string `yaml:"pid"`
	Auth     bool   `yaml:"auth"`
	Addr     string `yaml:"addr"`
	Mysql    string `yaml:"mysql"`
	LogDir   string `yaml:"log_dir"`
	LogLevel string `yaml:"log_level"`
}

var (
	confile = flag.String("c", "./app.conf.yaml", "配置文件路径")

	ClientConfig ClientConfigStruct
	ServeConfig  ServeConfigStruct
	Logger       *zap.Logger
	DB           *sqlx.DB
)

func init() {
	if e := gocommon.LoadYamlConfig(*confile, &ServeConfig); e != nil {
		panic(e)
	}

	if ServeConfig.LogDir != "" {
		if e := InitLog(); e != nil {
			panic(e)
		}
	}

	if ServeConfig.Mysql != "" {
		if e := InitDB(); e != nil {
			panic(e)
		}
	}

	if ServeConfig.Auth {
		passportcommon.Logger = Logger
		passportcommon.DB = DB
	}

	return
}

func InitLog() error {
	writer, _ := rotatelogs.New(
		ServeConfig.LogDir+"log.%Y%m%d%H%M",
		rotatelogs.WithLinkName("log"),
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	level := zapcore.DebugLevel
	if e := level.UnmarshalText([]byte(ServeConfig.LogLevel)); e != nil {
		return e
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(writer),
		level)

	Logger = zap.New(core, zap.Development())

	return nil
}

func InitDB() error {
	var e error

	if DB, e = sqlx.Connect("mysql", ServeConfig.Mysql); e != nil {
		return e
	}
	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	if e = DB.Ping(); e != nil {
		return e
	}

	return nil
}

func InitClient() error {
	if e := gocommon.LoadYamlConfig(*confile, &ClientConfig); e != nil {
		return e
	}

	writer, _ := rotatelogs.New(
		ServeConfig.LogDir+"log.%Y%m%d%H%M",
		rotatelogs.WithLinkName("log"),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(writer),
		zap.DebugLevel)

	Logger = zap.New(core, zap.Development())

	return nil
}

type MyLogger struct{}

func (p *MyLogger) Print(v ...interface{}) {
	Logger.Sugar().Debug(v)
}
