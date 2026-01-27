package logger

import (
	"os"
	"trae-go/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var L *zap.Logger

func InitLogger(cfg config.LogConfig, mode string) {
	// 1. 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 2. 配置 Encoder (日志格式)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式: 2024-01-01T12:00:00.000Z
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 级别格式: INFO, ERROR

	var encoder zapcore.Encoder
	if mode == "debug" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 开发模式用 Console 格式，好看点
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // 生产模式用 JSON 格式
	}

	// 3. 配置日志输出 (同时输出到文件和控制台)
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	})

	// 同时输出到控制台和文件
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)),
		level,
	)

	// 4. 创建 Logger
	// AddCaller: 添加调用者信息 (文件名:行号)
	L = zap.New(core, zap.AddCaller())

	// 替换全局的 logger，这样你也可以在其他地方用 zap.L() 直接调用
	zap.ReplaceGlobals(L)
}
