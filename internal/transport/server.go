package transport

import (
	"context"
	"github.com/segmentio/kafka-go"
	"net"
)

type Server struct {
	Handler *Handler
}

func NewServer(
	sessions SessionManager,
	producer *kafka.Writer,
	consumer *kafka.Reader,
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
