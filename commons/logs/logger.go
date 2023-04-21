package logs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

type Level uint8

const (
	LevelFatal Level = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var (
	levelMap = map[string]Level{
		"fatal": LevelFatal,
		"error": LevelError,
		"warn":  LevelWarning,
		"info":  LevelInfo,
		"debug": LevelDebug,
	}
	levelConvMap = map[Level]string{
		LevelFatal:   "fatal",
		LevelError:   "error",
		LevelWarning: "warn",
		LevelInfo:    "info",
		LevelDebug:   "debug",
	}
)

func (l Level) String() string {
	level := levelConvMap[l]
	if level != "" {
		return level
	}
	return "unknown"
}

type logger struct {
	level  Level
	log    func(calldepth int, s string) error
	caller int

	withPrefix bool
}

func NewDefaultLogger(withPrefix bool, level Level, log func(calldepth int, s string) error) Logger {
	r := &logger{withPrefix: withPrefix, level: level}
	r.log = log
	return r
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

var (
	defaultLogger = &logger{
		level:  LevelDebug,
		log:    log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Lmicroseconds).Output,
		caller: 4,
	}
)

func SetLevel(l Level) {
	defaultLogger.level = l
}
func SetLevelString(s string) {
	level, isExist := levelMap[strings.ToLower(s)]
	if !isExist {
		return
	}
	defaultLogger.level = level
}
func GetLevel() Level {
	return defaultLogger.level
}

func IsLevel(level Level) bool {
	return defaultLogger.level >= level
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}
func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}
func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxDebugf(ctx, format, v...)
}
func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxInfof(ctx, format, v...)
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.CtxErrorf(ctx, format, v...)
}

func (s *logger) Debugf(format string, v ...interface{}) {
	s.output(LevelDebug, fmt.Sprintf(format, v...))
}

func (s *logger) Infof(format string, v ...interface{}) {
	s.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (s *logger) Warnf(format string, v ...interface{}) {
	s.output(LevelWarning, fmt.Sprintf(format, v...))
}

func (s *logger) Errorf(format string, v ...interface{}) {
	s.output(LevelError, fmt.Sprintf(format, v...))
}

func (s *logger) Fatalf(format string, v ...interface{}) {
	s.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (s *logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	s.output(LevelDebug, fmt.Sprintf(GetLogId(ctx)+" "+format, v...))
}

func (s *logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	s.output(LevelInfo, fmt.Sprintf(GetLogId(ctx)+" "+format, v...))
}

func (s *logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	s.output(LevelError, fmt.Sprintf(GetLogId(ctx)+" "+format, v...))
}

func (s *logger) output(level Level, str string) {
	if level > s.level {
		return
	}
	formatStr := ""
	if s.withPrefix {
		switch level {
		case LevelFatal:
			formatStr = "\033[35m[FATAL]\033[0m " + str
		case LevelError:
			formatStr = "\033[31m[ERROR]\033[0m " + str
		case LevelWarning:
			formatStr = "\033[33m[WARN]\033[0m " + str
		case LevelInfo:
			formatStr = "\033[32m[INFO]\033[0m " + str
		case LevelDebug:
			formatStr = "\033[36m[DEBUG]\033[0m " + str
		}
	} else {
		formatStr = str
	}
	_ = s.log(s.caller, formatStr)
}

// log id

type logId struct {
}

var (
	traceId logId
)

func SetLogId(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceId, id)
}
func GetLogId(ctx context.Context) string {
	if ctx == nil {
		return "-"
	}
	result, _ := ctx.Value(traceId).(string)
	if result == "" {
		return "-"
	}
	return result
}
