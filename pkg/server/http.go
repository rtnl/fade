package server

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/fasthttp/websocket"
	"github.com/rtnl/fade/pkg/proto"
	"github.com/samber/mo"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{}

func (s *ServerImpl) RunHttp(ctx context.Context) mo.Result[any] {
	var (
		err error
	)

	err = fasthttp.ListenAndServe(":3000", s.handleHttpRequest)
	if err != nil {
		return mo.Err[any](err)
	}

	return mo.Ok[any](nil)
}

func (s *ServerImpl) handleHttpRequest(ctx *fasthttp.RequestCtx) {
	var (
		path string
	)

	path = string(ctx.Path())

	ctx.Response.Header.SetServer("fade")

	switch path {
	case "/":
		{
			s.handleHttpRequestIndex(ctx)
			break
		}

	case "/e":
		{
			s.handleHttpRequestExecute(ctx)
			break
		}

	case "/s":
		{
			s.handleHttpRequestSession(ctx)
			break
		}

	default:
		{
			ctx.Error("", 404)
			break
		}
	}
}

func (s *ServerImpl) handleHttpRequestIndex(ctx *fasthttp.RequestCtx) {
	ctx.Success("", []byte{})
}

func (s *ServerImpl) handleHttpRequestExecute(ctx *fasthttp.RequestCtx) {
	var (
		req       *proto.Req
		res       *proto.Res
		resFuture *mo.Future[*proto.Res]
		data      []byte
		err       error
	)

	data = ctx.Request.Body()

	err = json.Unmarshal(data, &req)
	if err != nil {
		slog.Info("failed at decoding req", slog.Any("err", err.Error()))
		ctx.Error("", 400)
		return
	}

	resFuture = s.executor.Handle(req)

	res, err = resFuture.Collect()
	if err != nil {
		slog.Info("failed at executing req", slog.Any("err", err.Error()))
		ctx.Error("", 400)
		return
	}

	data, err = json.Marshal(res)
	if err != nil {
		slog.Info("failed at encoding res", slog.Any("err", err.Error()))
		ctx.Error("", 400)
		return
	}

	slog.Info("successfully executed req", slog.Any("req_key", req.Key))
	ctx.Success("application/json", data)
	return
}

func (s *ServerImpl) handleHttpRequestSession(ctx *fasthttp.RequestCtx) {
	upgrader.Upgrade(ctx, s.handleWebsocketSession)
}
