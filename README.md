# paul-bunyan

Custom Echo API Logging Package

---

## Overview

The `bunyan` package is a customizable logger built around [logrus](https://github.com/sirupsen/logrus) and is developed to work with an [Echo](https://echo.labstack.com/) framework by replacing the default logger. 


This package includes:

1. A customizable logger to write messages
2. Middleware to log requests and responses


### Logrus
Logrus is a structured logger for Go. [logrus.Logger](https://github.com/sirupsen/logrus/blob/master/logger.go) 
is a struct that can be configured to set what, how, and where to write logs.

```golang
// Logrus Logger
type Logger struct {
	// Where to write the logs
	Out io.Writer

	// Hooks for the logger instance.
    // These allow firing events based on logging levels and log entries
	Hooks LevelHooks

	// How to write the logs
	// Included formatters are `TextFormatter` and `JSONFormatter`
	Formatter Formatter

	// Flag for whether to log caller info
	ReportCaller bool

	// The logging level the logger should log at
	Level Level

	// Used to sync writing to the log. Locking is enabled by Default
	mu MutexWrap

	// Reusable empty entry
	entryPool sync.Pool

	// Function to exit the application, defaults to `os.Exit()`
	ExitFunc exitFunc

	// The buffer pool used to format the log
	BufferPool BufferPool
}
```

---

## Customizing the Logger
Calling `bunyan.New()` returns a pointer to the APILogger struct which is the sole logger provided in this package. By default, logrus's logger is thread safe and there will be no race conditions. However, any formatting changes to this logger will persist and affect how all logs are written.

```golang
// Instantiate Logger
logger := bunyan.New()
```

### Levels

The five logging levels are as follows:

1. DEBUG
2. INFO
3. WARN
4. ERROR
5. OFF

The logger will only write logs of a level less than or equal to the set level. I.e., if the logger is set to `ERROR`, only logs at the `ERROR` level will be written and logs of levels `DEBUG`, `INFO`, or `WARN` will be exclude. Setting the level to `OFF` essentially mutes the logger.

To set the logger's level use the `SetLevel()` function.

```golang
logger.SetLevel(log.DEBUG)
logger.SetLevel(log.INFO)
logger.SetLevel(log.WARN)
logger.SetLevel(log.ERROR)
logger.SetLevel(log.OFF)
```

### Formats
Logrus has two built-in formatters.

1. [JSONFormatter](https://github.com/sirupsen/logrus/blob/master/json_formatter.go)
2. [TextFormatter](https://github.com/sirupsen/logrus/blob/master/text_formatter.go)

Custom formatters can also be created to further customize how the logs are written.

Use the `SetFormatter()` function to choose how the logs should look

```golang
logger.SetFormatter(&logrus.JSONFormatter{
    DisableHTMLEscape: true,
})
```

### Traceback information
To include the file name, line number, and function name as fields in the log, enable the traceback option.

```golang
logger.EnableTraceback()

// to turn off traceback
// logger.DisableTraceback()
```

### Output location
Use `SetOutput()` to write logs to a single location or `SetOutputs()` to write to multiple locations. Additionally, one can write to two or more locations using `SetOutput()` and passing an `io.MultiWriter`.

Typical options include:

1. os.Stdout
2. os.Stderr
3. file location using [lumberjack](https://github.com/natefinch/lumberjack)

```golang
// write to single location - standard error
logger.SetOutput(os.Stderr)

// write to 2 locations - file and standard out
logger.SetOutputs(&lumberjack.Logger{
    Filename:   logPath,
    MaxSize:    1, // megabytes
    MaxBackups: 100,
    MaxAge:     90,   //days
    Compress:   true, // disabled by default
}, os.Stdout)
```
### Lumberjack
[Lumberjack](https://github.com/natefinch/lumberjack) is a Go package for writing logs to rolling files. [lumberjack.Logger](https://github.com/natefinch/lumberjack/blob/47ffae23317c5951a2a6267a069cf676edf53eb6/lumberjack.go#L79) is a struct that configures the log file. Essentially, it will write to a specific file location and once that file reaches a certain size, it will start writing to a new file while retaining the old file for a set number of days.

```golang
// Lumberjack Logger
type Logger struct {
	// File to write logs to
	Filename string `json:"filename" yaml:"filename"`

	// Max size in megabytes of the log file before it gets rotated
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// determines if the time used in backup file names is the computer's local time
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// determines if the rotated log files should be compressed using gzip
	Compress bool `json:"compress" yaml:"compress"`

	size int64
	file *os.File
	mu   sync.Mutex

	millCh    chan bool
	startMill sync.Once
}
```

---

## Overwritting the default Echo Logger
With Echo, I recommend hidding the banner and port as these will also write to your logs on start and could make parsing logs difficult.

```golang
import (
	log "github.com/Dewberry/paul-bunyan"

	"github.com/labstack/echo/v4"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func main() {
    // Instantiate echo server
    e := echo.New()
    e.HideBanner = true
    e.HidePort = true

    // Instantiate Logger
    logger := log.New()

    // Set log level
    logger.SetLevel(log.DEBUG)

    // Set log output
    const logPath = "/logs/go-api.log"
    logger.SetOutputs(&lumberjack.Logger{
	    Filename:   logPath,
	    MaxSize:    1, // megabytes
	    MaxBackups: 100,
	    MaxAge:     90,   //days
	    Compress:   true, // disabled by default
    }, os.Stdout)

    // Will include file, line, and function fields with logs
    logger.EnableTraceback()

    // Set JSON formatter
    logger.SetFormatter(&logrus.JSONFormatter{
	    DisableHTMLEscape: true,
    })

    // Overwrite default logger
    e.Logger = logger
}
```

### Writing Logs within Handler Functions

```golang
func SomeHandler() echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Logger().Debug("Write a debug message here.")
        c.Logger().Info("Write an info message here.")
        c.Logger().Error("Write an error message here.")
        
        message := "Wrote some logs"
        return c.JSON(http.StatusOK, message)
    }
}
```

### Writing Logs in other functions
```golang
import myLogger "github.com/Dewberry/paul-bunyan"

func SomeHandler() echo.HandlerFunc {
    return func(c echo.Context) error {
        message := SomeFunction()
        
        return c.JSON(http.StatusOK, message)
    }
}

func SomeFunction() string {
    log := myLogger.New()

    log.Debug("Write a debug message here.")
    log.Info("Write an info message here.")
    log.Error("Write an error message here.")

    return "Wrote some logs"
}

```

---

## Logger Middleware
The logger middleware is used to write log both requests and responses. The middleware uses the same logger as mentioned above so any configurations applied (e.g., set output, level, format, etc.) will also be applied to these logs, traceback information being the exception.

There are two ways to use the logging middleware:

1. log.Middleware()
2. log.MiddlewareWithConfig(...)

```golang
import (
	log "github.com/Dewberry/paul-bunyan"

	"github.com/labstack/echo/v4"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func main() {
    // Instantiate echo server
    e := echo.New()
    
    // default middleware
    e.Use(log.Middleware())

    // configured middleware
    e.Use(
        log.MiddlewareWithConfig(
            // Request Config
            log.ReqConfig{
                Fields: []string {
                    "id",
                    "remote_ip",
                    "host",
                    "method",
                    "uri",
                    "user_agent",
                    "bytes_in",
                },
                Level: log.INFO,
                Message: "REQUEST",
            },
            // Response Config
            log.ResConfig{
                Fields: []string {
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
                Level: log.INFO,
                Message: "RESPONSE",
            },
        ),
    )
}
```

### Configuring the middleware

#### Fields -> fields to include in the logs

| Field | Description | Response | Request |
| ----- | ---------------------- | -------- | ------- |
|   id  | Request ID from header | <center>X</center> | <center>X</center> |
| remote_ip | Client's network address | <center>X</center> | <center>X</center> |
| host |URL's host name | <center>X</center> | <center>X</center> |
| method | HTTP method | <center>X</center> | <center>X</center> |
| uri | Request sent from client to server | <center>X</center> | <center>X</center> |
| user_agent | Client's user agent | <center>X</center> | <center>X</center> |
| status | HTTP status code |  | <center>X</center> |
| error | Any returned errors |  | <center>X</center> |
| latency | Latency in miliseconds |  | <center>X</center> |
| latency_human | Human-readable latency |  | <center>X</center> |
| bytes_in | Request size | <center>X</center> | <center>X</center> |
| bytes_out | Response size | | <center>X</center> |

By default, all fields above are included.

#### Level -> level to write logs

Set to DEBUG, INFO, TRACE, or ERROR to write the request or response log at that level. To mute either log set to OFF or set to a lower level than the logger's level (e.g., set request log's to write DEBUG and have the loggers level set to INFO).

By default, both are INFO level.

#### Message -> messages to include with the logs
Use to set the message field for request and response logs. 

By default the messages are "REQUEST" and "RESPONSE"