package ctx

import "context"

type RootContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

var rootCtx *RootContext = nil

func InitRootCtx() {
	if rootCtx != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	rootCtx = &RootContext{ctx, cancel}
}

func GetRootCtx() *RootContext {
	if rootCtx == nil {
		InitRootCtx()
	}

	return rootCtx
}
