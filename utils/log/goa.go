package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goadesign/goa"
)

type adapter struct {
	keyvals  []interface{}
	ctx      context.Context
	fallback goa.LogAdapter
}

func NewLogger(fallback goa.LogAdapter) goa.LogAdapter {
	if fallback == nil {
		fallback = goa.NewLogger(log.New(os.Stderr, "", log.LstdFlags))
	}
	return &adapter{
		fallback: fallback,
	}
}

type mainLog func(format string, data ...interface{})
type fallback func(msg string, keyvals ...interface{})

func (a *adapter) logMsg(whichFunc mainLog, fallbackFunc fallback, format string, data ...interface{}) {
	if a.ctx != nil {
		whichFunc(format, data...)
	} else {
		fallbackFunc(fmt.Sprintf(format, data...))
	}
}

func (a *adapter) debugf(format string, data ...interface{}) {
	a.logMsg(log.Printf, a.fallback.Info, format, data...)
}

func (a *adapter) infof(format string, data ...interface{}) {
	a.logMsg(log.Printf, a.fallback.Info, format, data...)
}

func (a *adapter) warningf(format string, data ...interface{}) {
	a.logMsg(log.Printf, a.fallback.Error, format, data...)
}

func (a *adapter) errorf(format string, data ...interface{}) {
	a.logMsg(log.Printf, a.fallback.Error, format, data...)
}

func (a *adapter) criticalf(format string, data ...interface{}) {
	a.logMsg(log.Printf, a.fallback.Error, format, data...)
}

func (a *adapter) Info(msg string, keyvals ...interface{}) {
	if a.ctx == nil {
		return
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, goa.ErrMissingLogValue)
	}
	if a.ctx != nil {
		a.infof("%s"+strings.Repeat(" %s=%+v", (len(a.keyvals)+len(keyvals))/2)+"\n", append([]interface{}{msg}, append(a.keyvals, keyvals...)...)...)
	} else {
		a.fallback.Info(msg, keyvals...)
	}
}

func (a *adapter) Error(msg string, keyvals ...interface{}) {
	if a.ctx == nil {
		return
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, goa.ErrMissingLogValue)
	}
	if a.ctx != nil {
		a.errorf("%s"+strings.Repeat(" %s=%+v", (len(a.keyvals)+len(keyvals))/2)+"\n", append([]interface{}{msg}, append(a.keyvals, keyvals...)...)...)
	} else {
		a.fallback.Error(msg, keyvals...)
	}
}

func (a *adapter) SetContext(ctx context.Context) *adapter {
	kvs := make([]interface{}, len(a.keyvals))
	copy(kvs, a.keyvals)
	return &adapter{
		keyvals:  kvs,
		ctx:      ctx,
		fallback: a.fallback,
	}
}

func (a *adapter) New(keyvals ...interface{}) goa.LogAdapter {
	if len(keyvals) == 0 {
		return a
	}
	kvs := append(a.keyvals, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, goa.ErrMissingLogValue)
	}

	return &adapter{
		keyvals:  kvs,
		ctx:      a.ctx,
		fallback: a.fallback.New(keyvals...),
	}
}
