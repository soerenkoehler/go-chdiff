package util

type ChainContext struct {
	Err error
}

func (ctx *ChainContext) Chain(f func()) *ChainContext {
	if ctx.Err == nil {
		f()
	}
	return ctx
}

func (ctx *ChainContext) ChainError(label string) {
	if ctx.Err != nil {
		Error("%s: %s", label, ctx.Err.Error())
	}
}

func (ctx *ChainContext) ChainFatal(label string) {
	if ctx.Err != nil {
		Fatal("%s: %s", label, ctx.Err.Error())
	}
}
