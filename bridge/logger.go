package logr

import (
	"github.com/echocat/slf4g"
	logr2 "github.com/echocat/slf4g-logr"
	"github.com/echocat/slf4g/level"
	"github.com/go-logr/logr"
)

// Default is the default instance which created with the slf4g global logger.
var Default = CreateFor(log.GetRootLogger())

func CreateFor(logger log.Logger) *Bridge {
	return &Bridge{
		Target: logger,
		VL:     0,
	}
}

type Bridge struct {
	Target      log.Logger
	VL          level.Level
	CallerDepth int
}

func (instance *Bridge) Enabled() bool {
	return instance.logger().IsLevelEnabled(instance.level() + level.Info)
}

func (instance *Bridge) log(l level.Level, msg string, err error, keysAndValues ...interface{}) {
	logger := instance.logger()
	f := logr2.KeysAndValuesToFields(keysAndValues...)
	if msg != "" {
		f[logger.GetProvider().GetFieldKeysSpec().GetMessage()] = msg
	}
	if err != nil {
		f[logger.GetProvider().GetFieldKeysSpec().GetError()] = err
	}
	logger.Log(logger.NewEvent(l, f), uint16(2+instance.CallerDepth))
}

func (instance *Bridge) Info(msg string, keysAndValues ...interface{}) {
	instance.log(instance.level()+level.Info, msg, nil, keysAndValues...)
}

func (instance *Bridge) Error(err error, msg string, keysAndValues ...interface{}) {
	instance.log(level.Error, msg, err, keysAndValues...)
}

func (instance *Bridge) V(l int) logr.Logger {
	return &Bridge{
		Target: instance.logger(),
		VL:     instance.VL + level.Level(l),
	}
}

func (instance *Bridge) WithValues(keysAndValues ...interface{}) logr.Logger {
	return &Bridge{
		Target: instance.logger().WithAll(logr2.KeysAndValuesToFields(keysAndValues...)),
		VL:     instance.VL,
	}
}

func (instance *Bridge) WithName(name string) logr.Logger {
	return &Bridge{
		Target: instance.logger().GetProvider().GetLogger(name),
		VL:     instance.VL,
	}
}

func (instance *Bridge) level() level.Level {
	if instance == nil {
		return 0
	}
	return instance.VL
}

func (instance *Bridge) logger() log.Logger {
	if instance == nil {
		return log.GetRootLogger()
	}
	if v := instance.Target; v != nil {
		return v
	}
	return log.GetRootLogger()
}
