package logger

import (
	"context"
	"os"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type LogHandler string

const (
	LogrusLogHandler LogHandler = "logrus"
)

type Fields map[string]interface{}
type LogLevel uint32

const (
	PanicLevel LogLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

const LogKey = "platform-logger"
const DefaultTimeStampFormat = time.RFC3339

func DefaultLoggerOption() NewLoggerOption {
	return NewLoggerOption{
		TimestampFormat: DefaultTimeStampFormat,
		LogHandler:      LogrusLogHandler,
	}
}

// Logger represent common interface for logging function
type ILogger interface {
	// Log a message at level Debug with optional additional fields.
	// This equivalent to `logger.WithField(fields).Info(msg)`
	Debug(msg string, field ...Fields)

	// Log a message at level Info with optional additional fields.
	// This equivalent to `logger.WithField(fields).Info(msg)`
	Info(msg string, field ...Fields)

	// Log a message at level Warn with optional additional fields.
	// This equivalent to `logger.WithField(fields).Warn(msg)`
	Warn(msg string, field ...Fields)

	// Log a message at level Error with optional additional fields.
	// This equivalent to `logger.WithField(fields).Error(msg)`
	Error(msg string, field ...Fields)
}

type Logger interface {
	ILogger

	// Context returns the context of the logger.
	Context() context.Context

	// WithContext returns a new Logger with the same handler
	// as the receiver and the given context.
	WithContext(ctx context.Context) Logger

	// AddLog Same as AddField but with the key is current time.
	// This suitable for specific tracing.
	// Example:
	// 	logger.AddLog("Step 01", nil)
	// 	logger.AddLog("Step 02",  User{
	//  		Name: "John Doe",
	//  		Age:  30,
	//  	})
	// The output as JSON:
	// 	{
	// 	  "2023-12-27 10:40:35.348336": "Step 01",
	// 	  "2023-12-27 10:40:35.348337": {
	// 	    "Step 02": {
	// 	      "Name": "John Doe",
	// 	      "Age": 30
	// 	    },
	// 	  }
	// 	}
	AddLog(msg string, value interface{}) Logger

	// AddField Add key/value field to log: key = value
	// Example:
	//
	// *With simple key/value*:
	//	logger.AddField("key", "value")
	//
	// *More complex key/value*:
	// type User struct {
	// 	Name string
	// 	Age  int
	// }
	// logger.AddField("key", User{
	// 	Name: "John Doe",
	// 	Age:  30,
	// })
	AddField(key string, value interface{}) Logger

	// AddFields Add multiple key/value to log: key1 = value1, key2 = value2
	// Example:
	//	logger.AddFields(plog.Fields{
	// 		"key1": "value1",
	// 		"key2": "value2",
	// 	})
	AddFields(fields Fields) Logger

	// WithField creates an new logger from the handler and adds a field to it.
	// If you want multiple fields, use `WithFields`.
	WithField(key string, value interface{}) Logger

	// WithField creates an new logger from the handler and adds multiple fields to it.
	// This is simply a helper for `WithField`, invoking it once for each field.
	WithFields(fields Fields) Logger
}

func NewLogger(opt NewLoggerOption) Logger {
	if opt.TimestampFormat == "" {
		opt.TimestampFormat = DefaultTimeStampFormat
	}
	if opt.LogHandler == "" {
		opt.LogHandler = LogrusLogHandler
	}
	switch opt.LogHandler {
	case LogrusLogHandler:
		logger := logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint:     false,
			TimestampFormat: opt.TimestampFormat,
		})
		entry := logger.WithFields(logrus.Fields(opt.DefaultFields))
		return &LogrusHandler{
			timestampFormat: opt.TimestampFormat,
			defaultEntry:    entry,
			Entry:           entry.Dup(),
		}
	default:
		panic("Unknown log handler: " + string(opt.LogHandler))
	}
}

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(NewLogger(DefaultLoggerOption()))
}

// Default returns the default Logger.
func Default() *LogrusHandler {
	return defaultLogger.Load().(*LogrusHandler)
}

type NewLoggerOption struct {
	TimestampFormat string
	DefaultFields   Fields
	LogHandler      LogHandler
}

// NewContext returns a context that contains the given Logger.
// Use FromContext to retrieve the Logger.
func NewContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, LogKey, l)
}

// FromContext returns the Logger stored in ctx by NewContext, or the default
// Logger if there is none.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(LogKey).(*LogrusHandler); ok {
		return l
	}
	return Default()
}

func GetLogger(ctx context.Context) Logger {
	return FromContext(ctx)
}
