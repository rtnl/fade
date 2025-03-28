package server

import (
	"context"
	"github.com/rtnl/fade/pkg/proto"
	"github.com/samber/mo"
)

type Handler interface {
	Init() mo.Result[any]
	Run(ctx context.Context, req *proto.Req) mo.Result[*proto.Res]
	GetKeys() []string
}
