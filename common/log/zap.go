package log

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"dmicro/common/config"
	"dmicro/common/config/env"
)

// 默认 debug level
var defaultLoggerConfig = config.Logger{
	Level:      "debug",
	Filename:   "app.log",
	MaxSize:    100,
	MaxBackups: 7,
	MaxAge:     7,
	Compress:   false,
}

var logger *zap.SugaredLogger

// Log level
func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Init(conf *config.Logger) {
	if conf == nil {
		conf = &defaultLoggerConfig
	}

	cores := []zapcore.Core{}

	{
		//if env.DeployEnv != env.DeployEnvProd {
		w := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			w,
			getLevel(conf.Level),
		)
		cores = append(cores, core)
		//}
	}

	{
		if env.DeployEnv == env.DeployEnvProd && conf.Filename != "" {
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   conf.Filename,
				MaxSize:    conf.MaxSize,
				MaxBackups: conf.MaxBackups,
				MaxAge:     conf.MaxAge,
				Compress:   conf.Compress,
			})
			core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				w,
				getLevel(conf.Level),
			)
			cores = append(cores, core)
		}
	}

	l := zap.New(zapcore.NewTee(cores...),
		zap.AddCallerSkip(1),
		zap.AddCaller(),
	)

	// defer with recover(): panic error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic recovered: %v\n", r)
			}
		}()
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			l.Sync()
		}
	}()

	logger = l.Sugar()
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
