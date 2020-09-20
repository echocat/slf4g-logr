package logr

import (
	"github.com/echocat/slf4g"
	logr2 "github.com/echocat/slf4g-logr"
	"github.com/go-logr/logr"
)

// Default is the default instance which created with the slf4g global logger.
var Default = CreateFor(log.GetGlobalLogger())

func CreateFor(logger log.Logger) *Bridge {
	return &Bridge{
		Target: logger,
		Level:  0,
	}
}

type Bridge struct {
	Target log.Logger
	Level  log.Level
}

func (instance *Bridge) Enabled() bool {
	return instance.logger().IsLevelEnabled(instance.level() + log.LevelInfo)
}

func (instance *Bridge) log(level log.Level, msg string, err error, keysAndValues ...interface{}) {
	logger := instance.logger()
	f := logr2.KeysAndValuesToFields(keysAndValues...)
	if msg != "" {
		f = f.With(logger.GetProvider().GetFieldKeySpec().GetMessage(), msg)
	}
	if err != nil {
		f = f.With(logger.GetProvider().GetFieldKeySpec().GetError(), err)
	}
	logger.Log(log.NewEvent(level, f, 3))
}

func (instance *Bridge) Info(msg string, keysAndValues ...interface{}) {
	instance.log(instance.level()+log.LevelInfo, msg, nil, keysAndValues...)
}

func (instance *Bridge) Error(err error, msg string, keysAndValues ...interface{}) {
	instance.log(log.LevelError, msg, err, keysAndValues...)
}

func (instance *Bridge) V(level int) logr.Logger {
	return &Bridge{
		Target: instance.logger(),
		Level:  instance.Level + log.Level(level),
	}
}

func (instance *Bridge) WithValues(keysAndValues ...interface{}) logr.Logger {
	return &Bridge{
		Target: instance.logger().WithFields(logr2.KeysAndValuesToFields(keysAndValues...)),
		Level:  instance.Level,
	}
}

func (instance *Bridge) WithName(name string) logr.Logger {
	return &Bridge{
		Target: instance.logger().GetProvider().GetLogger(name),
		Level:  instance.Level,
	}
}

func (instance *Bridge) level() log.Level {
	if instance == nil {
		return log.LevelInfo
	}
	if v := instance.Level; v != 0 {
		return v
	}
	return log.LevelInfo
}

func (instance *Bridge) logger() log.Logger {
	if instance == nil {
		return log.GetGlobalLogger()
	}
	if v := instance.Target; v != nil {
		return v
	}
	return log.GetGlobalLogger()
}
