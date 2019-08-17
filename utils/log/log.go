package log

import (
	"context"
	baseLog "log"

	"github.com/goadesign/goa"
)

func Debug(ctx context.Context, format string, data ...interface{}) {
	if l, ok := goa.ContextLogger(ctx).(*adapter); ok {
		l.debugf(format, data...)
	} else {
		baseLog.Println("Invalid logger for message:", format)
	}
}

func Info(ctx context.Context, format string, data ...interface{}) {
	if l, ok := goa.ContextLogger(ctx).(*adapter); ok {
		l.infof(format, data...)
	} else {
		baseLog.Println("Invalid logger for message:", format)
	}
}

func Warning(ctx context.Context, format string, data ...interface{}) {
	if l, ok := goa.ContextLogger(ctx).(*adapter); ok {
		l.warningf(format, data...)
	} else {
		baseLog.Println("Invalid logger for message:", format)
	}
}

func Error(ctx context.Context, format string, data ...interface{}) {
	if l, ok := goa.ContextLogger(ctx).(*adapter); ok {
		l.errorf(format, data...)
	} else {
		baseLog.Println("Invalid logger for message:", format)
	}
}

func Critical(ctx context.Context, format string, data ...interface{}) {
	if l, ok := goa.ContextLogger(ctx).(*adapter); ok {
		l.criticalf(format, data...)
	} else {
		baseLog.Println("Invalid logger for message:", format)
	}
}
