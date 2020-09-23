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
		VL:     0,
	}
}

type Bridge struct {
	Target      log.Logger
	VL          log.Level
	CallerDepth int
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
	logger.Log(log.NewEvent(level, f, 2+instance.CallerDepth))
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
		VL:     instance.VL + log.Level(level),
	}
}

func (instance *Bridge) WithValues(keysAndValues ...interface{}) logr.Logger {
	return &Bridge{
		Target: instance.logger().WithFields(logr2.KeysAndValuesToFields(keysAndValues...)),
		VL:     instance.VL,
	}
}

func (instance *Bridge) WithName(name string) logr.Logger {
	return &Bridge{
		Target: instance.logger().GetProvider().GetLogger(name),
		VL:     instance.VL,
	}
}

func (instance *Bridge) level() log.Level {
	if instance == nil {
		return 0
	}
	return instance.VL
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
