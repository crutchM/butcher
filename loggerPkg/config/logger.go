package config

import (
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"os"
)

const defaultLevel = slog.LevelDebug

var (
	ErrBadLogFormat = errors.New("неверное имя формата логирования")
	ErrBadLogLevel  = errors.New("неверный уровень логирования")
	defaultHandler  = slog.NewTextHandler(os.Stdout, nil)
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

const (
	JsonFormat = "json"
	TextFormat = "text"
)

type LoggerConfig struct {
	LogLevel string `yaml:"level"`
	Format   string `yaml:"format"`
}

func (lc LoggerConfig) GetLevel() (slog.Level, error) {
	switch lc.LogLevel {
	case LevelInfo:
		return slog.LevelInfo, nil

	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelError:
		return slog.LevelError, nil
	case LevelWarn:
		return slog.LevelWarn, nil

	default:
		return defaultLevel, fmt.Errorf("%v: получено значение - %s", ErrBadLogLevel, lc.LogLevel)
	}
}

func (lc LoggerConfig) GetHandler(opts *slog.HandlerOptions) (slog.Handler, error) {

	switch lc.Format {
	case JsonFormat:
		return slog.NewJSONHandler(os.Stdout, opts), nil

	case TextFormat:
		return slog.NewTextHandler(os.Stdout, opts), nil

	default:
		return defaultHandler, fmt.Errorf("%v: получено значение - %s", ErrBadLogFormat, lc.Format)
	}
}
