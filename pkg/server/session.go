package server

import (
	"context"
	"log/slog"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/rtnl/fade/pkg/proto"
)

type Session interface {
	GetId() uuid.UUID
	GetIdString() string
	GetConn() *websocket.Conn
	PushReq(value *proto.Req)
	PushRes(value *proto.Res)

	Init()
	Run()
	Wait()
	Stop()
}

type SessionImpl struct {
	executor *Executor

	ctx     context.Context
	cancel  context.CancelFunc
	id      uuid.UUID
	conn    *websocket.Conn
	chanReq chan *proto.Req
	chanRes chan *proto.Res
}

func NewSession(executor *Executor, conn *websocket.Conn) Session {
	s := new(SessionImpl)

	s.executor = executor

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.id = uuid.New()
	s.conn = conn
	s.chanReq = make(chan *proto.Req, 1024)
	s.chanRes = make(chan *proto.Res, 1024)

	return s
}

func (s *SessionImpl) GetId() uuid.UUID {
	return s.id
}

func (s *SessionImpl) GetIdString() string {
	return s.GetId().String()
}

func (s *SessionImpl) GetConn() *websocket.Conn {
	return s.conn
}

func (s *SessionImpl) PushReq(value *proto.Req) {
	if value == nil {
		return
	}

	s.chanReq <- value
}

func (s *SessionImpl) PushRes(value *proto.Res) {
	if value == nil {
		return
	}

	s.chanRes <- value
}

func (s *SessionImpl) Init() {
	slog.Info("created session",
		slog.Any("session_id", s.GetIdString()))
}

func (s *SessionImpl) runCyclePull() {
	var (
		req *proto.Req
		err error
	)

	defer s.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return

		default:
			break
		}

		req = nil

		err = s.GetConn().ReadJSON(&req)
		if err != nil {
			return
		}

		s.PushReq(req.Clone())
	}
}

func (s *SessionImpl) runCyclePush() {
	var (
		err error
	)

	defer s.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return

		case res := <-s.chanRes:
			{
				err = s.GetConn().WriteJSON(res)
				if err != nil {
					return
				}

				break
			}
		}
	}
}

func (s *SessionImpl) runCycleUpdate() {
	var (
		req *proto.Req
		res *proto.Res
		err error
	)

	defer s.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return

		case req = <-s.chanReq:
			{
				res, err = s.executor.Handle(req).Collect()
				if err != nil {
					return
				}

				s.PushRes(res)

				break
			}
		}
	}
}

func (s *SessionImpl) Run() {
	slog.Info("running session",
		slog.Any("session_id", s.GetIdString()))

	go s.runCyclePull()
	go s.runCyclePush()
	go s.runCycleUpdate()
}

func (s *SessionImpl) Wait() {
	<-s.ctx.Done()
}

func (s *SessionImpl) Stop() {
	if s.conn == nil {
		return
	}

	slog.Info("stopping session",
		slog.Any("session_id", s.GetIdString()))

	s.cancel()

	s.conn.Close()
	s.conn = nil
}
