package server

import (
	"context"
	"sync"

	"github.com/samber/mo"
)

type Server interface {
	Init() mo.Result[any]
	Run(ctx context.Context) mo.Result[any]

	AddHandler(value Handler)
	GetHandlerByKey(key string) mo.Option[Handler]
	ListHandler() []Handler
}

type ServerImpl struct {
	handlerList    []Handler
	handlerMap     map[string]Handler
	sessionMapLock sync.Mutex
	sessionMap     map[string]Session
	executor       *Executor
}

func NewServer() Server {
	s := &ServerImpl{
		handlerList:    make([]Handler, 0),
		handlerMap:     make(map[string]Handler),
		sessionMapLock: sync.Mutex{},
		sessionMap:     make(map[string]Session),
	}

	s.executor = NewExecutor(s)

	return s
}

func (s *ServerImpl) Init() mo.Result[any] {
	return mo.Ok[any](nil)
}

func (s *ServerImpl) Run(ctx context.Context) mo.Result[any] {
	go s.executor.Run(ctx)
	go s.RunHttp(ctx)

	<-ctx.Done()

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

func (s *ServerImpl) AddSession(value Session) {
	s.sessionMapLock.Lock()
	defer s.sessionMapLock.Unlock()

	s.sessionMap[value.GetIdString()] = value
}

func (s *ServerImpl) RemoveSession(key string) {
	s.sessionMapLock.Lock()
	defer s.sessionMapLock.Unlock()

	session, ok := s.sessionMap[key]
	if !ok {
		return
	}

	session.Stop()
	delete(s.sessionMap, key)
}
