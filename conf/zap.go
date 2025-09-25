package conf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type LogConfig struct {
	Stdout bool
	// DebugLevel, InfoLevel, WarnLevel, ErrorLevel, PanicLevel
	level            zapcore.Level
	caller           bool
	lumberjackConfig *LumberjackConfig
}

type LumberjackConfig struct {
	Filename string
	// rotate if current file too large
	MaxSizeInMB int
	// can keep how many old log files
	MaxBackups int
	// compress old log files
	Compress bool
}

var Stdout *zap.Logger = NewLog(LogConfig{
	Stdout:           true,
	level:            zap.DebugLevel,
	caller:           true,
	lumberjackConfig: nil,
})

var FileLog *zap.Logger = NewLog(LogConfig{
	Stdout: false,
	level:  zap.DebugLevel,
	caller: true,
	lumberjackConfig: &LumberjackConfig{
		Filename:    ".zap-log",
		MaxSizeInMB: 8,
		MaxBackups:  7,
		Compress:    false,
	},
})

func NewLog(cfg LogConfig) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.TimeKey = "time"
	ec.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) {
		e.AppendString(t.In(time.Local).Format("2006-01-02 15:04:05"))
	}
	ec.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewJSONEncoder(ec)
	if cfg.Stdout {
		encoder = zapcore.NewConsoleEncoder(ec)
	}
	var ws []zapcore.WriteSyncer
	if cfg.lumberjackConfig != nil {
		fileRotate := &lumberjack.Logger{
			Filename:   cfg.lumberjackConfig.Filename,
			MaxSize:    cfg.lumberjackConfig.MaxSizeInMB,
			MaxBackups: cfg.lumberjackConfig.MaxBackups,
			LocalTime:  true,
			Compress:   cfg.lumberjackConfig.Compress,
		}
		ws = append(ws, zapcore.AddSync(fileRotate))
	}
	if cfg.Stdout {
		ws = append(ws, zapcore.Lock(os.Stdout))
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(ws...), cfg.level)
	return zap.New(core, zap.WithCaller(cfg.caller))
}
