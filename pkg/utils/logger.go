/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 01:28
@Name: logger.go
@Description:
*/

package utils

import (
	"context"
	"github.com/openzipkin/zipkin-go/idgenerator"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
)

var (
	LogName    = "server.log"
	LogPath    = "log"
	defaultLog = GetLogger(context.Background())
)

func GetLogger(ctx context.Context, taskID ...string) logrus.FieldLogger {
	if len(taskID) > 0 {
		logger := NewLog(ctx, LogPath, taskID[0])
		ctx = context.WithValue(ctx, TraceID, taskID[0])
		ctx = context.WithValue(ctx, Logger, logger)
		return logger
	}
	return NewLog(ctx, LogPath, DefaultLogger)
}

func GetDefaultLogger() logrus.FieldLogger {
	return defaultLog
}

func NewLog(ctx context.Context, logPath string, value string) logrus.FieldLogger {
	// TODO: 滚动日志,实现按天切割日志文件
	// https://github.com/natefinch/lumberjack
	// https://github.com/lestrrat-go/file-rotatelogs

	loggerInterface := ctx.Value(value)
	if loggerInterface != nil {
		return loggerInterface.(logrus.FieldLogger)
	}
	traceID := idgenerator.NewRandomTimestamped().TraceID().String()
	logger := InitLogger(logPath, traceID)
	ctx = context.WithValue(ctx, TraceID, traceID)
	ctx = context.WithValue(ctx, value, logger)
	return logger
}

func InitLogger(logPath string, traceID string) logrus.FieldLogger {
	lg := logrus.New()
	lg.SetOutput(io.MultiWriter(os.Stdout,
		&lumberjack.Logger{
			Filename:   filepath.Join(logPath, LogName),
			MaxSize:    10,
			MaxAge:     3,
			MaxBackups: 3,
			LocalTime:  false,
			Compress:   true,
		}))

	lg.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z",
		FullTimestamp:   true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyFile:  "file",
			logrus.FieldKeyMsg:   "msg",
			TraceID:              traceID,
		},
		DisableColors: false,
	})

	return lg
}
