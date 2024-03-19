package logger

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LogrusHandler struct {
	m               sync.Mutex
	timestampFormat string
	defaultEntry    *logrus.Entry
	Entry           *logrus.Entry
	ctx             context.Context
}

func (l *LogrusHandler) GetTimestampFormat() string {
	return l.timestampFormat
}

func (l *LogrusHandler) AddLog(msg string, value interface{}) Logger {
	const TimeFormat = "2006-01-02 15:04:05.000000"
	logTime := time.Now().Format(TimeFormat)
	if value == nil {
		l.AddField(logTime, msg)
	} else {
		l.AddField(logTime, Fields{msg: value})
	}
	return l
}

func (l *LogrusHandler) AddField(key string, value interface{}) Logger {
	l.m.Lock()
	defer l.m.Unlock()
	l.Entry = l.Entry.WithField(key, value)
	return l
}

// Add multiple key/value to log: key1 = value1, key2 = value2
func (l *LogrusHandler) AddFields(fields Fields) Logger {
	l.m.Lock()
	defer l.m.Unlock()
	l.Entry = l.Entry.WithFields(logrus.Fields(fields))
	return l
}

func (l *LogrusHandler) WithField(key string, value interface{}) Logger {
	return l.withEntry(l.Entry.WithField(key, value))
}

func (l *LogrusHandler) WithFields(fields Fields) Logger {
	return l.withEntry(l.Entry.WithFields(logrus.Fields(fields)))
}

func (l *LogrusHandler) Debug(msg string, field ...Fields) {
	l.log(msg, logrus.DebugLevel, field...)
}

func (l *LogrusHandler) Info(msg string, field ...Fields) {
	l.log(msg, logrus.InfoLevel, field...)
}

func (l *LogrusHandler) Warn(msg string, field ...Fields) {
	l.log(msg, logrus.WarnLevel, field...)
}

func (l *LogrusHandler) Error(msg string, field ...Fields) {
	l.log(msg, logrus.ErrorLevel, field...)
}

func (l *LogrusHandler) log(msg string, level logrus.Level, field ...Fields) {
	entry := l.Entry.Dup()
	for _, f := range field {
		entry = entry.WithFields(logrus.Fields(f))
	}

	entry.Log(logrus.Level(level), msg)
}

func (l *LogrusHandler) withEntry(entry *logrus.Entry) *LogrusHandler {
	return &LogrusHandler{
		timestampFormat: l.timestampFormat,
		defaultEntry:    l.defaultEntry,
		Entry:           entry,
	}
}

func (l *LogrusHandler) dup() *LogrusHandler {
	return &LogrusHandler{
		timestampFormat: l.timestampFormat,
		defaultEntry:    l.defaultEntry,
		Entry:           l.Entry.Dup(),
		ctx:             l.ctx,
	}
}

func (l *LogrusHandler) WithContext(ctx context.Context) Logger {
	dup := l.dup()
	dup.ctx = ctx
	return dup
}

func (l *LogrusHandler) Context() context.Context {
	return l.ctx
}
