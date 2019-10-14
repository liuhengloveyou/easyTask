package common

import (
	"flag"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	gocommon "github.com/liuhengloveyou/go-common"
	passportcommon "github.com/liuhengloveyou/passport/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type ServeConfigStruct struct {
	PID       string `yaml:"pid"`
	Addr      string `yaml:"addr"`
	Mysql     string `yaml:"mysql"`
	LogDir    string `yaml:"log_dir"`
	LogLevel  string `yaml:"log_level"`
}

type ClientConfigStruct struct {
	PID       string `yaml:"pid"`
	Host      string `yaml:"host"`
	Addr      string `yaml:"addr"`
	Mysql     string `yaml:"mysql"`
	LogDir    string `yaml:"log_dir"`
	LogLevel  string `yaml:"log_level"`
}

var (
	confile = flag.String("c", "./app.conf.yaml", "配置文件路径")

	ServeConfig ServeConfigStruct
	Logger     *zap.Logger
)

func init() {
	if e := gocommon.LoadYamlConfig(*confile, &ServeConfig); e != nil {
		panic(e)
	}

	writer, _ := rotatelogs.New(
		ServeConfig.LogDir+"app.log.%Y%m%d%H%M",
		rotatelogs.WithLinkName("app.log"),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(writer),
		zap.DebugLevel)

	Logger = zap.New(core, zap.Development())

	passportcommon.Logger = Logger

	return
}

type MyLogger struct {}
func (p *MyLogger) Print(v ...interface{}) {
	Logger.Sugar().Debug(v)
}
