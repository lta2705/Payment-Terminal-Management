package transport

import (
	"Payment-Terminal-Management/internal/service"
	"context"
	"github.com/Jeffail/gabs"
	"github.com/bytedance/gopkg/util/logger"
	"net"
)

type SessionManager interface {
	Add(conn net.Conn, terminalId string)
	Remove(key string)
	Check(terminalId string) bool
	Send(trmId string, message string) bool
	Count() int
	Broadcast(message string)
	CloseAll()
}

type Handler struct {
	Sessions      SessionManager
	KafkaProducer service.ProducerService
	KafkaConsumer service.ConsumerService
}

func NewHandler(
	sessions SessionManager,
	kafkaProducer service.ProducerService,
	kafkaConsumer service.ConsumerService,
) *Handler {
	return &Handler{
		Sessions:      sessions,
		KafkaProducer: kafkaProducer,
		KafkaConsumer: kafkaConsumer,
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

			h.processMessage(conn, buf[:n])
		}
	}
}

func (h *Handler) processMessage(conn net.Conn, data []byte) {
	logger.Info(string(data))

	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		logger.Error("Invalid message format", err)
		return
	}

	msgTypeRaw := jsonParsed.Search("msgType").Data()
	msgType, ok := msgTypeRaw.(string)
	if !ok {
		logger.Error("msgType is missing or not a string")
		return
	}

	switch msgType {
	case "0":
		trmIdRaw := jsonParsed.Path("trmId").Data()
		trmId, ok := trmIdRaw.(string)
		if !ok || trmId == "" {
			logger.Error("TerminalId field missing or not a string")
			return
		}

		h.Sessions.Add(conn, trmId)
		logger.Info("Session opened:", trmId)

	case "2":
		err := h.KafkaProducer.ProduceMessage(string(data))
		if err != nil {
			logger.Error("Failed to send message to Kafka", err)
		}

	default:
		logger.Warn("Unknown MsgType:", msgType)
	}

}
func (h *Handler) Close() {
	h.Sessions.CloseAll()
}
