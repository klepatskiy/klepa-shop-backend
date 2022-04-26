package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
)

const logPath string = "APP_LOG_PATH"

type zapLogger struct {
	loggers map[string]*zap.Logger
}

type ZapLogger interface {
	Error(loggerName string, msg string, fields ...zap.Field)
	Info(loggerName string, msg string, fields ...zap.Field)
	Fatal(loggerName string, msg string, fields ...zap.Field)
}

func NewZapLogger() ZapLogger {
	loggers := make(map[string]*zap.Logger)
	return &zapLogger{loggers: loggers}
}

func (zl *zapLogger) Error(loggerName string, msg string, fields ...zap.Field) {
	zl.getZapLogger(loggerName).Error(msg, fields...)
}

func (zl *zapLogger) Info(loggerName string, msg string, fields ...zap.Field) {
	zl.getZapLogger(loggerName).Info(msg, fields...)
}

func (zl *zapLogger) Fatal(loggerName string, msg string, fields ...zap.Field) {
	zl.getZapLogger(loggerName).Fatal(msg, fields...)
}

func (zl *zapLogger) getZapLogger(loggerName string) *zap.Logger {
	zapLogger, ok := zl.loggers[loggerName]
	if ok == true {
		return zapLogger
	}

	cfg := zap.NewProductionConfig()

	cfg.Encoding = "console"
	//cfg.DisableStacktrace = true

	path := getEnv(logPath, nil)

	if path != nil {
		strPath := fmt.Sprintf("%v", path) + loggerName + ".Log"

		_, _ = CreatFileWithFolders(strPath)

		cfg.OutputPaths = []string{
			strPath,
			"stdout",
		}
	} else {
		cfg.OutputPaths = []string{
			"stdout",
		}
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build()

	if err == nil && path == nil {
		logger.Error(
			"Env parameter for Log file not set",
			zap.String("env_parameter", logPath),
		)
	}

	if err != nil {
		log.Fatal("Error loading zapLogger %V",
			err,
		)
	}

	zl.loggers[loggerName] = logger

	return logger
}

func getEnv(key string, defaultVal interface{}) interface{} {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func CreatFileWithFolders(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil { //nolint:gosec
		return nil, err
	}
	file, _ := os.Open(p) //nolint:gosec
	if file != nil {
		return file, nil
	}
	return os.Create(p) //nolint:gosec
}
