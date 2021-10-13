package trace

//trace for labstack/echo

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func EchoTraceMiddleWare(key string) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.New().String()
			enty := logrus.WithFields(logrus.Fields{
				TRACE_ID: id,
			})
			c.Set(TRACE_ID, id)
			c.Set(TRACE_LOGGER, enty)
			return h(c)
		}
	}
}

func GetLogger(c echo.Context) *logrus.Entry {
	l, ok := c.Get(TRACE_LOGGER).(*logrus.Entry)
	if !ok {
		_, l = setNewTrace(c)
	}
	return l
}

func GetTraceID(c echo.Context) string {
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
