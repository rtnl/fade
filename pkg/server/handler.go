package server

import (
	"github.com/samber/mo"
)

type Handler interface {
	Init() mo.Result[any]
	Run() mo.Result[any]
	GetKeys() []string
}
