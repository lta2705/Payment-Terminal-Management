package transport

import (
	"Payment-Terminal-Management/internal/dto"
	"context"
	"encoding/json"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
	"net"
)

type SessionManager interface {
	Add(conn net.Conn, terminalId string)
	Remove(key string)
	CloseAll()
}

type Handler struct {
	Sessions      SessionManager
	KafkaProducer *kafka.Writer
	KafkaConsumer *kafka.Reader
}

func NewHandler(
	sessions SessionManager,
	producer *kafka.Writer,
	consumer *kafka.Reader,
) *Handler {
	return &Handler{
		Sessions:      sessions,
		KafkaProducer: producer,
		KafkaConsumer: consumer,
	}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	addr := conn.RemoteAddr().String()

	defer func() {
		h.Sessions.Remove(addr)
		conn.Close()
		logger.Info("Session closed:", addr)
	}()

	buf := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				return
			}

			h.processMessage(conn, addr, buf[:n])
		}
	}
}

func (h *Handler) processMessage(conn net.Conn, addr string, data []byte) {
	var msg dto.RegisterDTO

	if err := json.Unmarshal(data, &msg); err != nil {
		logger.Error("Invalid message format", err)
		return
	}

	switch msg.MsgType {
	case "0":
		h.Sessions.Add(conn, msg.TerminalId)
		logger.Info("Session opened:", msg.TerminalId)

	case "2":

	default:
		logger.Warn("Unknown MsgType:", msg.MsgType)
	}

	conn.Write([]byte("ACK\n"))
}

func (h *Handler) Close() {
	h.Sessions.CloseAll()
}
