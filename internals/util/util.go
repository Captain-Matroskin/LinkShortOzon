package util

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	errPkg "linkShortOzon/internals/myerror"
)

const (
	LenLinkShort = 10
	LinkDomain   = "ozon.click.ru" // test domain
)

type Logger struct {
	Log errPkg.MultiLoggerInterface
}

func NewLogger(filePath string) *zap.SugaredLogger {
	configLog := zap.NewProductionEncoderConfig()
	configLog.TimeKey = "time_stamp"
	configLog.LevelKey = "level"
	configLog.MessageKey = "note"
	configLog.EncodeTime = zapcore.ISO8601TimeEncoder
	configLog.EncodeLevel = zapcore.CapitalLevelEncoder

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     60,
		Compress:   false,
	}
	writerSyncer := zapcore.AddSync(lumberJackLogger)
	encoder := zapcore.NewConsoleEncoder(configLog)

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	zapLogger := logger.Sugar()
	return zapLogger
}
