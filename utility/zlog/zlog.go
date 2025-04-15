package zlog

import (
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func write_syncer(fn string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename: fn,
		MaxSize:  100, //MB default value is also 100
		//MaxBackups: 100,
		MaxAge:    90, //days
		LocalTime: true,
		Compress:  true,
	})
}

func Init(fn string, no_std, no_trace bool) {
	enc_cfg := zapcore.EncoderConfig{
		// some of the following fileds are used only by json encoder
		TimeKey:          "tm",
		LevelKey:         "lvl",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stack",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
	}

	// file output
	enc := zapcore.NewConsoleEncoder(enc_cfg)

	default_level := zap.InfoLevel
	if !no_std {
		default_level = zap.DebugLevel
	}
	info_c := zapcore.NewCore(enc, write_syncer("./log/info/"+fn+".log"), default_level)
	err_c := zapcore.NewCore(enc, write_syncer("./log/err/"+fn+".err.log"), zap.ErrorLevel)

	core := zapcore.NewTee(info_c, err_c)

	if !no_std {
		enc_cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		std_enc := zapcore.NewConsoleEncoder(enc_cfg)
		std_syncer := zapcore.AddSync(os.Stdout)
		std_c := zapcore.NewCore(std_enc, std_syncer, zap.DebugLevel)
		core = zapcore.NewTee(core, std_c)

		// enc_cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		// debug_c := zapcore.NewCore(std_enc, write_syncer("./log/debug/"+fn+".debug.log"), zap.DebugLevel)
		// core = zapcore.NewTee(core, debug_c)
	}

	if no_trace {
		logger = zap.New(core)
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1) /*, zap.AddStacktrace(zap.ErrorLevel)*/)
	}
	sugar = logger.Sugar()
}

func Sync() {
	logger.Sync()
	sugar.Sync()
}

func Info(message string, fields ...zap.Field) {
	// callers := trace()
	// fields = append(fields, callers...)
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	// callers := trace()
	// fields = append(fields, callers...)
	logger.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	// callers := trace()
	// fields = append(fields, callers...)
	logger.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	// callers := trace()
	// fields = append(fields, callers...)
	logger.Warn(message, fields...)
}

// !!! panic() will be called after the log output
func Panic(message string, fields ...zap.Field) {
	// callers := trace()
	// fields = append(fields, callers...)
	logger.Panic(message, fields...)
}

func Infof(format string, v ...any) {
	sugar.Infof(format, v...)
}

func Debugf(format string, v ...any) {
	sugar.Debugf(format, v...)
}

func Errorf(format string, v ...any) {
	sugar.Errorf(format, v...)
}

func Warnf(format string, v ...any) {
	sugar.Warnf(format, v...)
}

func trace() (fields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2) //skip 2 levels
	if !ok {
		return
	}
	func_path := runtime.FuncForPC(pc).Name()
	func_name := path.Base(func_path) //Base return the last elememnt of path which is also func name

	fields = append(fields, zap.String("func", func_name), zap.String("file", file), zap.Int("line", line))
	return
}
