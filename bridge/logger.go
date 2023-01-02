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
	VL          int
	CallerDepth int
}

func (instance *Bridge) Init(info logr.RuntimeInfo) {
	instance.CallerDepth = info.CallDepth
}

func (instance *Bridge) Enabled(vl int) bool {
	return instance.logger().IsLevelEnabled(instance.level(vl))
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

func (instance *Bridge) Info(lvl int, msg string, keysAndValues ...interface{}) {
	instance.log(instance.level(lvl), msg, nil, keysAndValues...)
}

func (instance *Bridge) Error(err error, msg string, keysAndValues ...interface{}) {
	instance.log(level.Error, msg, err, keysAndValues...)
}

func (instance *Bridge) V(vl int) logr.LogSink {
	return &Bridge{
		Target: instance.logger(),
		VL:     instance.vl(vl),
	}
}

func (instance *Bridge) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &Bridge{
		Target: instance.logger().WithAll(logr2.KeysAndValuesToFields(keysAndValues...)),
		VL:     instance.vl(0),
	}
}

func (instance *Bridge) WithName(name string) logr.LogSink {
	return &Bridge{
		Target: instance.logger().GetProvider().GetLogger(name),
		VL:     instance.vl(0),
	}
}

func (instance *Bridge) vl(vl int) int {
	if instance == nil {
		return vl
	}
	return instance.VL + vl
}

func (instance *Bridge) level(vl int) level.Level {
	vl = instance.vl(vl)
	if vl < 0 {
		return level.Warn
	}

	switch vl {
	case 0, 1:
		return level.Info
	case 2, 3:
		return level.Debug
	case 4, 5:
		return level.Trace
	default:
		return 0
	}

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
