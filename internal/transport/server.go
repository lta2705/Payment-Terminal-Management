package transport

import (
	"Payment-Terminal-Management/internal/service"
	"context"
	"net"
)

type Server struct {
	Handler *Handler
}

func NewServer(
	sessions SessionManager,
	producer service.ProducerService,
	consumer service.ConsumerService,
) *Server {
	return &Server{
		Handler: NewHandler(sessions, producer, consumer),
	}
}

func (s *Server) HandleConnection(ctx context.Context, conn net.Conn) {
	s.Handler.Handle(ctx, conn)
}

func (s *Server) Close() {
	s.Handler.Close()
}
