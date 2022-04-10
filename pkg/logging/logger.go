package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger = nil

type mode int8

const (
	DEBUG mode = iota
	PROD
)

func createPath(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return err
	}
	_, err = os.Create(path)
	return err
}

func setTimeEncoding(config *zap.Config, m mode) {
	switch m {
	case PROD:

		config.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
			// 2019-08-13T04:39:11Z
		})
	case DEBUG:
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.Stamp))
			// Aug 13 00:38:11
		})
	}
}

func InitLoggerConfig(jsonConfig []byte) {
	var cfg zap.Config
	if err := json.Unmarshal(jsonConfig, &cfg); err != nil {
		panic(err)
	}

	if err := createPath(cfg.OutputPaths[1]); err != nil {
		panic(err)
	}

	setTimeEncoding(&cfg, mode(0))
	loggerRaw, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = loggerRaw.Sugar()
}

func InitDefaultLogger() {
	loggerRaw, error := zap.NewDevelopment()
	if error != nil {
		fmt.Println(error)
	}
	//defer logger.Sync()
	logger = loggerRaw.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		InitDefaultLogger()
	}
	return logger
}
