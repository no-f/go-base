package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/no-f/go-base/config/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
)

var Logger *zap.Logger
var level int

// Initialize 日志初始化
func Initialize(loggerConfig *models.LoggerYAMLConfig) {
	logDir := loggerConfig.LogPath
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	logFilePath := filepath.Join(logDir, loggerConfig.LogName)
	level = loggerConfig.LogLevel

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 控制台日志核心
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		getLogLevel(),
	)

	// 文件日志核心
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    loggerConfig.LogMaxSize,
			MaxBackups: loggerConfig.LogMaxBackups,
			MaxAge:     loggerConfig.LogMaxAge,
			Compress:   loggerConfig.LogCompress,
		}),
		getLogLevel(),
	)

	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	if err := logger.Sync(); err != nil {
		log.Fatalf("同步日志失败: %v", err)
	}

	Logger = logger
}

func getLogLevel() zap.AtomicLevel {
	return zap.NewAtomicLevelAt(GetZapLogLevel())
}

func GetZapLogLevel() zapcore.Level {
	switch level {
	case -1:
		return zapcore.DebugLevel
	case 0:
		return zapcore.InfoLevel
	case 1:
		return zapcore.WarnLevel
	case 2:
		return zapcore.ErrorLevel
	case 3:
		return zapcore.DPanicLevel
	case 4:
		return zapcore.PanicLevel
	case 5:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Info(title string, fields ...zap.Field) {
	Logger.Info(title, fields...)
}

func Warn(title string, fields ...zap.Field) {
	Logger.Warn(title, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
	os.Exit(1)
}
