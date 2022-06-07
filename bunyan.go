package bunyan

import (
	"encoding/json"
	"io"
	"runtime"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

const (
	DEBUG log.Lvl = iota + 1
	INFO
	WARN
	ERROR
	OFF
)

// select which fields to include
// select level
type (
	APILogger struct {
		*logrus.Logger
	}
	ReqConfig struct {
		Fields  []string
		Level   log.Lvl
		Message string
	}
	ResConfig struct {
		Fields  []string
		Level   log.Lvl
		Message string
	}
)

var (
	customLogger = &APILogger{
		Logger: logrus.New(),
	}
	DefaultReqConfig = ReqConfig{
		Fields: []string{
			"id",
			"remote_ip",
			"host",
			"method",
			"uri",
			"user_agent",
			"bytes_in",
		},
		Level:   INFO,
		Message: "REQUEST",
	}
	DefaultResConfig = ResConfig{
		Fields: []string{
			"id",
			"remote_ip",
			"host",
			"method",
			"uri",
			"user_agent",
			"status",
			"error",
			"latency",
			"latency_human",
			"bytes_in",
			"bytes_out",
		},
		Level:   INFO,
		Message: "RESPONSE",
	}
	traceback bool = false
)

func New() *APILogger {
	return customLogger
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Enables writing file, line, and function name with logs
func (l *APILogger) EnableTraceback() {
	traceback = true
}

// Disables writing file, line, and function name with logs
func (l *APILogger) DisableTraceback() {
	traceback = false
}

// Used to get file, line number, and func that sent log message
func getTraceback() logrus.Fields {
	// skip 2 levels to get file, line, and func that sent log
	// level 1 would return the log functions below
	if traceback {
		skip := 2

		pc, file, line, _ := runtime.Caller(skip)
		funcName := runtime.FuncForPC(pc).Name()

		return logrus.Fields{"file": file, "line": line, "func": funcName}
	}
	return logrus.Fields{}
}

// Print output message of print level
func Print(i ...interface{}) {
	customLogger.Print(i...)
}

// Printf output format message of print level
func Printf(format string, i ...interface{}) {
	customLogger.Printf(format, i...)
}

// Printj output json of print level
func Printj(j log.JSON) {
	customLogger.Printj(j)
}

// Debug output message of debug level
func Debug(i ...interface{}) {
	customLogger.Debug(i...)
}

// Debugf output format message of debug level
func Debugf(format string, args ...interface{}) {
	customLogger.Debugf(format, args...)
}

// Debugj output json of debug level
func Debugj(j log.JSON) {
	customLogger.Debugj(j)
}

// Info output message of info level
func Info(i ...interface{}) {
	customLogger.Info(i...)
}

// Infof output format message of info level
func Infof(format string, args ...interface{}) {
	customLogger.Infof(format, args...)
}

// Infoj output json of info level
func Infoj(j log.JSON) {
	customLogger.Infoj(j)
}

// Warn output message of warn level
func Warn(i ...interface{}) {
	customLogger.Warn(i...)
}

// Warnf output format message of warn level
func Warnf(format string, args ...interface{}) {
	customLogger.Warnf(format, args...)
}

// Warnj output json of warn level
func Warnj(j log.JSON) {
	customLogger.Warnj(j)
}

// Error output message of error level
func Error(i ...interface{}) {
	customLogger.Error(i...)
}

// Errorf output format message of error level
func Errorf(format string, args ...interface{}) {
	customLogger.Errorf(format, args...)
}

// Errorj output json of error level
func Errorj(j log.JSON) {
	customLogger.Errorj(j)
}

// Fatal output message of fatal level
func Fatal(i ...interface{}) {
	customLogger.Fatal(i...)
}

// Fatalf output format message of fatal level
func Fatalf(format string, args ...interface{}) {
	customLogger.Fatalf(format, args...)
}

// Fatalj output json of fatal level
func Fatalj(j log.JSON) {
	customLogger.Fatalj(j)
}

// Panic output message of panic level
func Panic(i ...interface{}) {
	customLogger.Panic(i...)
}

// Panicf output format message of panic level
func Panicf(format string, args ...interface{}) {
	customLogger.Panicf(format, args...)
}

// Panicj output json of panic level
func Panicj(j log.JSON) {
	customLogger.Panicj(j)
}

// To logrus.Level
func toLogrusLevel(level log.Lvl) logrus.Level {
	switch level {
	case log.DEBUG:
		return logrus.DebugLevel
	case log.INFO:
		return logrus.InfoLevel
	case log.WARN:
		return logrus.WarnLevel
	case log.ERROR:
		return logrus.ErrorLevel
	}

	return logrus.InfoLevel
}

// To Echo.log.lvl
func toEchoLevel(level logrus.Level) log.Lvl {
	switch level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	}

	return log.OFF
}

// Output return logger io.Writer
func (l *APILogger) Output() io.Writer {
	return l.Out
}

// SetOutput logger io.Writer
func (l *APILogger) SetOutput(w io.Writer) {
	l.Out = w
}

// SetOutput logger io.Writer
func (l *APILogger) SetOutputs(w ...io.Writer) {
	l.Out = io.MultiWriter(w...)
}

// Level return logger level
func (l *APILogger) Level() log.Lvl {
	return toEchoLevel(l.Logger.Level)
}

// SetLevel logger level
func (l *APILogger) SetLevel(v log.Lvl) {
	l.Logger.Level = toLogrusLevel(v)
}

// SetHeader logger header
// Managed by Logrus itself
// This function do nothing
func (l *APILogger) SetHeader(h string) {
	// do nothing
}

// Formatter return logger formatter
func (l *APILogger) Formatter() logrus.Formatter {
	return l.Logger.Formatter
}

// SetFormatter logger formatter
// Only support logrus formatter
func (l *APILogger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = formatter
}

// Prefix return logger prefix
// This function do nothing
func (l *APILogger) Prefix() string {
	return ""
}

// SetPrefix logger prefix
// This function do nothing
func (l *APILogger) SetPrefix(p string) {
	// do nothing
}

// Print output message of print level
func (l *APILogger) Print(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Print(i...)
}

// Printf output format message of print level
func (l *APILogger) Printf(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Printf(format, args...)
}

// Printj output json of print level
func (l *APILogger) Printj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Println(string(b))
}

// Debug output message of debug level
func (l *APILogger) Debug(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Debug(i...)
}

// Debugf output format message of debug level
func (l *APILogger) Debugf(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Debugf(format, args...)
}

// Debugj output message of debug level
func (l *APILogger) Debugj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Debugln(string(b))
}

// Info output message of info level
func (l *APILogger) Info(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Info(i...)
}

// Infof output format message of info level
func (l *APILogger) Infof(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Infof(format, args...)
}

// Infoj output json of info level
func (l *APILogger) Infoj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Infoln(string(b))
}

// Warn output message of warn level
func (l *APILogger) Warn(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Warn(i...)
}

// Warnf output format message of warn level
func (l *APILogger) Warnf(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Warnf(format, args...)
}

// Warnj output json of warn level
func (l *APILogger) Warnj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Warnln(string(b))
}

// Error output message of error level
func (l *APILogger) Error(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Error(i...)
}

// Errorf output format message of error level
func (l *APILogger) Errorf(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Errorf(format, args...)
}

// Errorj output json of error level
func (l *APILogger) Errorj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Errorln(string(b))
}

// Fatal output message of fatal level
func (l *APILogger) Fatal(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Fatal(i...)
}

// Fatalf output format message of fatal level
func (l *APILogger) Fatalf(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Fatalf(format, args...)
}

// Fatalj output json of fatal level
func (l *APILogger) Fatalj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Fatalln(string(b))
}

// Panic output message of panic level
func (l *APILogger) Panic(i ...interface{}) {
	l.Logger.WithFields(getTraceback()).Panic(i...)
}

// Panicf output format message of panic level
func (l *APILogger) Panicf(format string, args ...interface{}) {
	l.Logger.WithFields(getTraceback()).Panicf(format, args...)
}

// Panicj output json of panic level
func (l *APILogger) Panicj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.WithFields(getTraceback()).Panicln(string(b))
}

// Logger Middleware Function
func Middleware() echo.MiddlewareFunc {
	return MiddlewareWithConfig(DefaultReqConfig, DefaultResConfig)
}

// Logger Middleware Function
func MiddlewareWithConfig(reqConfig ReqConfig, resConfig ResConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			reqFields := logrus.Fields{}
			resFields := logrus.Fields{}

			// Get request and response message
			reqMessage := reqConfig.Message
			resMessage := resConfig.Message

			// Retrieve info to pass with log as fields

			// id field
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "id") {
				reqFields["id"] = id
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "id") {
				resFields["id"] = id
			}

			// remote_ip field
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "remote_ip") {
				reqFields["remote_ip"] = c.RealIP()
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "remote_ip") {
				resFields["remote_ip"] = c.RealIP()
			}

			// host field
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "host") {
				reqFields["host"] = req.Host
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "host") {
				resFields["host"] = req.Host
			}

			// bytes_in field
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "bytes_in") {
				reqFields["bytes_in"] = reqSize
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "bytes_in") {
				resFields["bytes_in"] = reqSize
			}

			// method field
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "method") {
				reqFields["method"] = req.Method
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "method") {
				resFields["method"] = req.Method
			}

			// uri field
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "uri") {
				reqFields["uri"] = req.RequestURI
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "uri") {
				resFields["uri"] = req.RequestURI
			}

			// user_agent field
			if reqConfig.Fields == nil || contains(reqConfig.Fields, "user_agent") {
				reqFields["user_agent"] = req.UserAgent()
			}
			if resConfig.Fields == nil || contains(resConfig.Fields, "user_agent") {
				resFields["user_agent"] = req.UserAgent()
			}

			// Send request log
			switch reqConfig.Level {
			case DEBUG:
				customLogger.WithFields(reqFields).Debug(reqMessage)
			case INFO:
				customLogger.WithFields(reqFields).Info(reqMessage)
			case WARN:
				customLogger.WithFields(reqFields).Warn(reqMessage)
			case ERROR:
				customLogger.WithFields(reqFields).Error(reqMessage)
			case OFF:
				// do nothing
			default:
				customLogger.WithFields(reqFields).Info(reqMessage)
			}

			start := time.Now()
			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			// status field
			if resConfig.Fields == nil || contains(resConfig.Fields, "status") {
				resFields["status"] = res.Status
			}

			// error field
			if resConfig.Fields == nil || contains(resConfig.Fields, "error") {
				switch err {
				case nil:
					resFields["error"] = ""
				default:
					resFields["error"] = err
				}
			}

			// latency field
			if resConfig.Fields == nil || contains(resConfig.Fields, "latency") {
				resFields["latency"] = stop.Sub(start)
			}

			// latency_human field
			if resConfig.Fields == nil || contains(resConfig.Fields, "latency_human") {
				resFields["latency_human"] = stop.Sub(start).String()
			}

			// bytes_out field
			if resConfig.Fields == nil || contains(resConfig.Fields, "bytes_out") {
				resFields["bytes_out"] = strconv.FormatInt(res.Size, 10)
			}

			// Send response log
			switch reqConfig.Level {
			case DEBUG:
				customLogger.WithFields(resFields).Debug(resMessage)
			case INFO:
				customLogger.WithFields(resFields).Info(resMessage)
			case WARN:
				customLogger.WithFields(resFields).Warn(resMessage)
			case ERROR:
				customLogger.WithFields(resFields).Error(resMessage)
			case OFF:
				// do nothing
			default:
				customLogger.WithFields(resFields).Info(resMessage)
			}

			return err
		}
	}
}
