package trace

//trace for labstack/echo

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func EchoTraceMiddleWare(key string) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := uuid.New()
			c.Set(TRACE_ID, id)
			return h(c)
		}
	}
}
