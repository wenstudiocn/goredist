package log

import (
	"github.com/wenstudiocn/goredist/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

const (
	DEF_LOG_FILENAME = "log.log"
	DEF_LOG_PATH     = "./"
	DEF_LOG_LEVEL    = 1
)

var (
	defaultLog *SLogger

	LogLevelMap = map[int]zapcore.Level{
		1: zapcore.DebugLevel,
		2: zapcore.InfoLevel,
		3: zapcore.WarnLevel,
		4: zapcore.ErrorLevel,
		5: zapcore.DPanicLevel,
		6: zapcore.PanicLevel,
		7: zapcore.FatalLevel,
	}
)

/// signle logger
/// means log in a single file

///  默认配置:
///  json 格式
///  单文件最大 10 m
/// 最多保留1个月
/// 压缩备份
type SLogger struct {
	f  string
	lg *zap.Logger
}

// @console: 是否输出到 console
// @path: 路径
// @level: 日志等级
// @sinks: 日志额外的输出
func NewSLogger(console bool, logfile string, level int, sinks ...zap.Sink) *SLogger {
	//修正参数
	// log file
	dir, filename := path.Split(logfile)
	if dir == "" {
		dir = DEF_LOG_PATH
	}
	if filename == "" {
		filename = DEF_LOG_FILENAME
	}
	err := utils.EnsurePath(dir)
	if err != nil {
		return nil
	}
	fullpath := path.Join(dir, filename)
	// log level
	if level < 1 || level > 7 {
		level = DEF_LOG_LEVEL
	}
	settings := lumberjack.Logger{
		Filename:   fullpath,
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 360,
		Compress:   true,
	}

	encoderConf := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// level
	atom := zap.NewAtomicLevelAt(zapcore.Level(level))
	wss := []zapcore.WriteSyncer{zapcore.AddSync(&settings)}
	if console {
		wss = append(wss, zapcore.AddSync(os.Stdout))
	}
	for _, sink := range sinks {
		wss = append(wss, zapcore.AddSync(sink))
	}
	// writeSyncer
	ws := zapcore.NewMultiWriteSyncer(wss...)
	// core
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf),
		ws,
		atom)
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(2)
	dev := zap.Development()
	stack := zap.AddStacktrace(zapcore.ErrorLevel)

	logger := zap.New(core, caller, callerSkip, dev, stack)

	return &SLogger{
		f:  fullpath,
		lg: logger,
	}
}

func (self *SLogger) Debug(msg string, fields ...zap.Field) {
	self.lg.Debug(msg, fields...)
}

func (self *SLogger) Info(msg string, fields ...zap.Field) {
	self.lg.Info(msg, fields...)
}

func (self *SLogger) Warn(msg string, fields ...zap.Field) {
	self.lg.Warn(msg, fields...)
}

func (self *SLogger) Error(msg string, fields ...zap.Field) {
	self.lg.Error(msg, fields...)
}

func (self *SLogger) Panic(msg string, fields ...zap.Field) {
	self.lg.Panic(msg, fields...)
}

func (self *SLogger) Fatal(msg string, fields ...zap.Field) {
	self.lg.Fatal(msg, fields...)
}

func (self *SLogger) Sync() {
	_ = self.lg.Sync()
}

func SetDefaultLogger(logger *SLogger) {
	defaultLog = logger
}

func Debug(msg string, fields ...zap.Field) {
	defaultLog.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLog.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLog.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLog.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	defaultLog.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLog.Fatal(msg, fields...)
}

func Sync() {
	defaultLog.Sync()
}
