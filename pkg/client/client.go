package client

import (
	"context"

	"github.com/rtnl/fade/pkg/proto"
	"github.com/samber/mo"
)

type ClientExecuteMode int32

const (
	ClientExecuteModeNone ClientExecuteMode = iota
	ClientExecuteModeHttp
	ClientExecuteModeHttps
	ClientExecuteModeWs
	ClientExecuteModeWss
)

var (
	ClientExecuteOrder = []ClientExecuteMode{
		ClientExecuteModeHttps,
	}
)

type ClientExecuteContext struct {
	req     *proto.Req
	attempt int
}

func NewClientExecuteContext(req *proto.Req) (ctx *ClientExecuteContext) {
	ctx = new(ClientExecuteContext)

	ctx.req = req
	ctx.attempt = 1

	return
}

type Client interface {
	Init() mo.Result[any]
	Run(ctx context.Context) mo.Result[any]
	executeContext(ctx *ClientExecuteContext) *mo.Future[*proto.Res]
	Execute(req *proto.Req) *mo.Future[*proto.Res]
	ExecuteMode(req *proto.Req, mode ClientExecuteMode) *mo.Future[*proto.Res]
	ExecuteForward(req *proto.Res) mo.Result[any]

	GetExecuteChannel() chan *proto.Res
}

type ClientImpl struct {
}

func NewClient() Client {
	c := new(ClientImpl)

	return c
}

func (c *ClientImpl) Init() mo.Result[any] {
	return mo.Ok[any](nil)
}

func (c *ClientImpl) Run(ctx context.Context) mo.Result[any] {
	return mo.Ok[any](nil)
}

func (c *ClientImpl)

func (c *ClientImpl) Execute(req *proto.Req) *mo.Future[*proto.Res] {
	return nil
}

func (c *ClientImpl) ExecuteMode(req *proto.Req, mode ClientExecuteMode) *mo.Future[*proto.Res] {
	return nil
}

func (c *ClientImpl) ExecuteForward(req *proto.Res) mo.Result[any] {
	return mo.Ok[any](nil)
}
