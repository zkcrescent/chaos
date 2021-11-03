package trace

import "context"

type Context struct {
	context.Context
}

func NewContext() CCTX {
	return &Context{
		context.Background(),
	}
}

func (c *Context) SetCTX(k, v interface{}) {
	c.Context = context.WithValue(c.Context, k, v)
}

func (c *Context) GetCTX(k interface{}) interface{} {
	return c.Context.Value(k)
}
