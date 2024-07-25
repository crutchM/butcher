package sentrylog

import (
	"butcher/loggerPkg"
	"butcher/loggerPkg/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type _sentryLog struct {
	// клиент сентри
	client *sentry.Client
	// откуда создали клиент(нужен для указания модуля)
	scope string
	ctx   context.Context
}

func (sentryLog _sentryLog) extractStructFromArgs(event *sentry.Event, args ...any) []any {
	var newArgs []any
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		if v.Kind() == reflect.Struct {
			if jsonData, err := json.Marshal(arg); err == nil {
				event.Extra[fmt.Sprintf("struct_%T", arg)] = string(jsonData)
			}
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs
}

func (sentryLog _sentryLog) getEvent(msg string, level sentry.Level, err error, args ...any) *sentry.Event {
	event := sentry.NewEvent()
	newArgs := sentryLog.extractStructFromArgs(event, args)
	event.Message = fmt.Sprintf("message: %s, args: %v", msg, newArgs)
	event.Level = level
	if err != nil {
		event.Exception = []sentry.Exception{
			{
				Type:   string(level),
				Value:  err.Error(),
				Module: sentryLog.scope,
			},
		}
	}
	event.EventID = sentry.EventID(uuid.New().String())

	return event
}

func (sentryLog _sentryLog) Info(msg string, args ...any) {
	event := sentryLog.getEvent(msg, sentry.LevelInfo, nil, args)
	sentryLog.client.CaptureEvent(event, nil, sentry.NewScope())
}

func (sentryLog _sentryLog) Debug(msg string, args ...any) {
	event := sentryLog.getEvent(msg, sentry.LevelDebug, nil, args)
	sentryLog.client.CaptureEvent(event, nil, sentry.NewScope())

}

func (sentryLog _sentryLog) Warn(msg string, err error, args ...any) {
	event := sentryLog.getEvent(msg, sentry.LevelWarning, err, args)
	sentryLog.client.CaptureEvent(event, nil, sentry.NewScope())

}

func (sentryLog _sentryLog) Error(msg string, err error, args ...any) {
	event := sentryLog.getEvent(msg, sentry.LevelError, err, args)
	sentryLog.client.CaptureEvent(event, nil, sentry.NewScope())

}

func NewSentryLog(ctx context.Context, cfg config.Sentry, host string, orgID int) (loggerPkg.Callback, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Debug:            cfg.Debug,
		Dsn:              cfg.DSN, // URL для подключения к сервису в сентри
		AttachStacktrace: true,
		ServerName:       host,            // хост с которого прилетит ивент в сам сентри, добавляем для идентификации нужного адаптера
		Environment:      cfg.Environment, // здесь указывается с какого окружения отправляем ивенты - прод/стейдж
	})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка создание клиента sentry логгера: ")
	}

	cx := sentry.SetHubOnContext(ctx, sentry.CurrentHub().Clone())

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("organization-id", strconv.Itoa(orgID))
	})

	return &_sentryLog{client: client, ctx: cx}, nil
}
