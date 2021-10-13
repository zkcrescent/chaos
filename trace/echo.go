package trace

//trace for labstack/echo

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type CCTX context.Context

type CTX struct {
	echo.Context
	CCTX
}

func EchoTraceMiddleWare() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.New().String()
			enty := logrus.WithFields(logrus.Fields{
				TRACE_ID: id,
			})
			c.Set(TRACE_ID, id)
			c.Set(TRACE_LOGGER, enty)
			c = &CTX{
				Context: c,
				CCTX:    context.Background(),
			}
			return h(c)
		}
	}
}

func SetLogger(c interface{}, l *logrus.Entry) {
	ctx := getCTXOrPanic(c)
	ctx.Set(TRACE_LOGGER, l)
}

func getCTXOrPanic(in interface{}) *CTX {
	return in.(*CTX)
}

func GetLogger(in interface{}) *logrus.Entry {
	ctx := getCTXOrPanic(in)

	l, ok := ctx.Get(TRACE_LOGGER).(*logrus.Entry)
	if !ok {
		_, l = setNewTrace(ctx)
	}
	return l
}

func GetTraceID(in interface{}) string {
	c := getCTXOrPanic(in)
	str, ok := c.Get(TRACE_ID).(string)
	if !ok {
		str, _ = setNewTrace(c)
	}
	return str

}

func setNewTrace(c echo.Context) (string, *logrus.Entry) {
	id := uuid.New().String()
	l := logrus.WithFields(logrus.Fields{
		TRACE_ID: id,
	})
	c.Set(TRACE_ID, id)
	c.Set(TRACE_LOGGER, l)
	return id, l
}

func Start(ctx interface{}) *CTX {
	t := getCTXOrPanic(ctx)
	return t
}
