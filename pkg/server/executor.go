package server

import (
	"context"
	"fmt"
	"time"

	"github.com/rtnl/fade/pkg/proto"
	"github.com/samber/mo"
)

type ExecutorEntry struct {
	req     *proto.Req
	resolve func(*proto.Res)
	reject  func(error)
}

func NewExecutorEntry(req *proto.Req, resolve func(*proto.Res), reject func(error)) *ExecutorEntry {
	return &ExecutorEntry{
		req:     req,
		resolve: resolve,
		reject:  reject,
	}
}

type Executor struct {
	server *ServerImpl

	queue chan *ExecutorEntry
}

func NewExecutor(server *ServerImpl) (e *Executor) {
	e = new(Executor)

	e.server = server
	e.queue = make(chan *ExecutorEntry)

	return
}

func (e *Executor) Init() mo.Result[any] {
	return mo.Ok[any](nil)
}

func (e *Executor) Run(ctx context.Context) {
	var (
		entry *ExecutorEntry
	)

	ticker := time.NewTicker(time.Millisecond * 1)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			break
		}

		entry = <-e.queue
		if entry == nil {
			continue
		}

		e.Execute(ctx, entry)
	}
}

func (e *Executor) Handle(req *proto.Req) *mo.Future[*proto.Res] {
	return mo.NewFuture[*proto.Res](func(resolve func(*proto.Res), reject func(error)) {
		e.queue <- NewExecutorEntry(req, resolve, reject)
	})
}

func (e *Executor) Execute(ctx context.Context, entry *ExecutorEntry) {
	var (
		req     *proto.Req
		res     *proto.Res
		handler Handler
		ok      bool
		err     error
	)

	req = entry.req

	handler, ok = e.server.GetHandlerByKey(req.GetMethod()).Get()
	if !ok {
		entry.reject(fmt.Errorf("handler not found"))
		return
	}

	res, err = handler.Run(ctx, req).Get()
	if err != nil {
		entry.reject(fmt.Errorf("handler error: %s", err.Error()))
		return
	}

	entry.resolve(res)
	return
}
