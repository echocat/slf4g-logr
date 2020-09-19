package logr

import (
	"github.com/echocat/slf4g"
	"github.com/go-logr/logr"
)

// Default is the default instance which created with the slf4g global logger.
var Default = CreateFor(log.GetGlobalLogger())

func CreateFor(logger log.Logger) *Bridge {
	return &Bridge{
		Logger: logger,
		Level:  log.LevelInfo,
	}
}

type Bridge struct {
	Logger log.Logger
	Level  log.Level
}

func (instance *Bridge) Enabled() bool {
	return instance.logger().IsLevelEnabled(instance.level())
}

func (instance *Bridge) Info(msg string, keysAndValues ...interface{}) {
	instance.Logger.
		WithFields(KeysAndValuesToFields(keysAndValues...)).
		Info(msg)
}

func (instance *Bridge) Error(err error, msg string, keysAndValues ...interface{}) {
	instance.Logger.
		WithError(err).
		WithFields(KeysAndValuesToFields(keysAndValues...)).
		Error(msg)
}

func (instance *Bridge) V(level int) logr.Logger {
	return &Bridge{
		Logger: instance.Logger,
		Level:  log.Level(level),
	}
}

func (instance *Bridge) WithValues(keysAndValues ...interface{}) logr.Logger {
	return &Bridge{
		Logger: instance.Logger.WithFields(KeysAndValuesToFields(keysAndValues...)),
		Level:  instance.Level,
	}
}

func (instance *Bridge) WithName(name string) logr.Logger {
	return &Bridge{
		Logger: instance.Logger.GetProvider().GetLogger(name),
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
	if v := instance.Logger; v != nil {
		return v
	}
	return log.GetGlobalLogger()
}
