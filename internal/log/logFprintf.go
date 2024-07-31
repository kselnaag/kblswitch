package log

import (
	"fmt"
	T "kblswitch/internal/types"
	"os"
	"time"
)

var _ T.ILog = (*LogFprintf)(nil)

type LogFprintf struct {
	loglvl  T.LogLevel
	appname string
}

func NewLogFprintf(appname string, loglvl string) *LogFprintf {
	var lvl T.LogLevel
	switch loglvl {
	case T.StrTrace:
		lvl = T.Trace
	case T.StrDebug:
		lvl = T.Debug
	case T.StrInfo:
		lvl = T.Info
	case T.StrWarn:
		lvl = T.Warn
	case T.StrError:
		lvl = T.Error
	case T.StrFatal:
		lvl = T.Fatal
	case T.StrPanic:
		lvl = T.Panic
	default:
		lvl = T.Logsoff
	}
	return &LogFprintf{
		loglvl:  lvl,
		appname: appname,
	}
}

func logMessage(lvl, svc, mess, err string) {
	timenow := time.Now().Format(time.RFC3339Nano)
	fmt.Fprintf(os.Stderr, `{"T":"%s","L":"%s","S":"%s","M":"%s","E":"%s"}`+"\n", timenow, lvl, svc, mess, err)
}

func (l *LogFprintf) LogTrace(format string, v ...any) {
	if l.loglvl <= T.Trace {
		logMessage(T.StrTrace, l.appname, fmt.Sprintf(format, v...), "")
	}
}

func (l *LogFprintf) LogDebug(format string, v ...any) {
	if l.loglvl <= T.Debug {
		logMessage(T.StrDebug, l.appname, fmt.Sprintf(format, v...), "")
	}
}

func (l *LogFprintf) LogInfo(format string, v ...any) {
	if l.loglvl <= T.Info {
		logMessage(T.StrInfo, l.appname, fmt.Sprintf(format, v...), "")
	}
}

func (l *LogFprintf) LogWarn(format string, v ...any) {
	if l.loglvl <= T.Warn {
		logMessage(T.StrWarn, l.appname, fmt.Sprintf(format, v...), "")
	}
}

func (l *LogFprintf) LogError(err error, format string, v ...any) {
	if l.loglvl <= T.Error {
		logMessage(T.StrError, l.appname, fmt.Sprintf(format, v...), err.Error())
	}
}

func (l *LogFprintf) LogFatal(err error, format string, v ...any) {
	if l.loglvl <= T.Fatal {
		logMessage(T.StrFatal, l.appname, fmt.Sprintf(format, v...), err.Error())
	}
}

func (l *LogFprintf) LogPanic(err error, format string, v ...any) {
	if l.loglvl <= T.Panic {
		logMessage(T.StrPanic, l.appname, fmt.Sprintf(format, v...), err.Error())
	}
}
