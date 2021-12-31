package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

func Init(config *Config) {
	logger = NewSugaredLogger(config)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

type Config struct {
	Sugar  bool
	Level  int
	File   string
	Path   string
	Rotate RotateConfig
}

type RotateConfig struct {
	Compress   bool
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

func NewLogger(conf *Config) *zap.Logger {
	return initZapLogger(conf)
}

func NewSugaredLogger(conf *Config) *zap.SugaredLogger {
	log := initZapLogger(conf)
	return log.Sugar()
}

func initZapLogger(conf *Config) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   conf.Path,
		MaxSize:    conf.Rotate.MaxSize,
		MaxBackups: conf.Rotate.MaxBackups,
		MaxAge:     conf.Rotate.MaxAge,
		Compress:   conf.Rotate.Compress,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(conf.Level))

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		atomicLevel,
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.Development())
	return zapLogger
}
