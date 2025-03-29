package server

import (
	"github.com/fasthttp/websocket"
)

func (s *ServerImpl) handleWebsocketSession(conn *websocket.Conn) {
	var (
		session Session
	)

	session = NewSession(s.executor, conn)
	session.Init()

	s.AddSession(session)

	session.Run()
	session.Wait()

	s.RemoveSession(session.GetIdString())
}
