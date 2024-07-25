package slogPkg

import (
	"butcher/loggerPkg"
	"butcher/loggerPkg/config"
	"github.com/pkg/errors"
	"log/slog"
)

type _slog struct {
	*slog.Logger
}

func (l _slog) Warn(msg string, err error, args ...any) {
	var log *slog.Logger

	if err != nil {
		log = l.With(slog.Any("error", err.Error()))
	} else {
		log = l.Logger
	}

	log.Warn(msg, args)
}

func (l _slog) Error(msg string, err error, args ...any) {
	var log *slog.Logger

	if err != nil {
		log = l.With(slog.Any("error", err.Error()))
	} else {
		log = l.Logger
	}

	log.Error(msg, args)
}

func NewSlog(cfg config.LoggerConfig) (loggerPkg.Callback, error) {
	l := new(_slog)

	level, err := cfg.GetLevel()

	if err != nil {
		return nil, errors.Wrap(err, "ошибка инициализации стандратного логгера")
	}

	handler, err := cfg.GetHandler(&slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
	l.Logger = slog.New(handler)

	return l, nil
}
