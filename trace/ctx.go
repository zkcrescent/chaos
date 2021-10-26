package trace

import "context"

type Context struct {
	ctx context.Context
}

func NewContext() CCTX {
	return &Context{
		context.Background(),
	}
}

func (c *Context) SetCTX(k, v interface{}) {
	context.WithValue(c.ctx, k, v)
}

func (c *Context) GetCTX(k interface{}) interface{} {
	return c.ctx.Value(k)
}
