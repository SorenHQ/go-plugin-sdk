package logtool

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ServiceName string
	DevMode     bool
	sugar       *zap.SugaredLogger
)

func GetLogger() *zap.SugaredLogger {
	return sugar
}

func Init(serviceName string, devMode bool) {
	ServiceName = serviceName
	DevMode = devMode
	var logger *zap.Logger
	var err error
	if devMode {
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.TimeKey = ""
		zapConfig.OutputPaths = []string{"stdout"}

		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		zapConfig.EncoderConfig.EncodeName = zapcore.FullNameEncoder
		zapConfig.DisableStacktrace = true
		logger, err = zapConfig.Build()
	} else {
		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "time"
		zapConfig.Encoding = "json"
		zapConfig.OutputPaths = []string{"stdout"}
		logger, err = zapConfig.Build()
	}
	if err != nil {
		log.Fatal(err)
	}
	sugar = logger.Sugar()
}

// Custom Fiber logger middleware for zap
