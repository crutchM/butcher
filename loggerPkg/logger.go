package loggerPkg

import (
	"butcher/loggerPkg/config"
	"butcher/loggerPkg/slogPkg"
)

type Callback interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, err error, args ...any)
	Error(msg string, err error, args ...any)
}

type Logger struct {
	context   []any
	callbacks []Callback
}

func NewLogger(cfg config.LoggerConfig) (*Logger, error) {
	logger := new(Logger)

	sl, err := slogPkg.NewSlog(cfg)

	if err != nil {
		return nil, err
	}
	logger.callbacks = append(logger.callbacks, sl)

	return logger, nil
}

func (log *Logger) LoggerWith(ctxArgs ...any) *Logger {

	return &Logger{
		context:   ctxArgs,
		callbacks: log.callbacks,
	}
}

func (log *Logger) LoggerWithCallBack(callbacks ...Callback) *Logger {
	var lw Logger
	lw.callbacks = make([]Callback, 0, len(callbacks))
	for _, callback := range callbacks {
		lw.callbacks = append(lw.callbacks, callback)
	}
	return &lw
}

func (log *Logger) Info(msg string, args ...any) {
	for _, callback := range log.callbacks {
		callback.Info(msg, args)
	}
}

func (log *Logger) Debug(msg string, args ...any) {
	//var attrs []any
	//if len(log.context) !=0{
	//	attrs = log.getArgsWithAttributes(args)
	//} else {
	//	attrs = args
	//}
	for _, callback := range log.callbacks {
		callback.Debug(msg, args)
	}
}

func (log *Logger) Warn(msg string, err error, args ...any) {
	for _, callback := range log.callbacks {
		callback.Warn(msg, err, args)
	}
}

func (log *Logger) Error(msg string, err error, args ...any) {

	for _, callback := range log.callbacks {
		callback.Error(msg, err, args)
	}
}
