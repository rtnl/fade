package server

import (
	"github.com/samber/mo"
)

type Server interface {
	Init() mo.Result[any]
	Run() mo.Result[any]

	AddHandler(value Handler)
	GetHandlerByKey(key string) mo.Option[Handler]
	ListHandler() []Handler
}

type ServerImpl struct {
	handlerList []Handler
	handlerMap  map[string]Handler
}

func NewServer() Server {
	return &ServerImpl{
		handlerList: make([]Handler, 0),
		handlerMap:  make(map[string]Handler),
	}
}

func (s *ServerImpl) Init() mo.Result[any] {
	return mo.Ok[any](nil)
}

func (s *ServerImpl) Run() mo.Result[any] {
	return mo.Ok[any](nil)
}

func (s *ServerImpl) AddHandler(value Handler) {
	if value == nil {
		return
	}

	s.handlerList = append(s.handlerList, value)

	for _, key := range value.GetKeys() {
		s.handlerMap[key] = value
	}
}

func (s *ServerImpl) GetHandlerByKey(key string) mo.Option[Handler] {
	value, ok := s.handlerMap[key]

	if !ok {
		return mo.None[Handler]()
	} else {
		return mo.Some(value)
	}
}

func (s *ServerImpl) ListHandler() []Handler {
	return s.handlerList
}
