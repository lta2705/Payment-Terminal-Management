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
	logger.Info("Object after parsed", jsonParsed)
	MsgType := jsonParsed.S("msg_type").Data().(string)

	switch MsgType {
	case "0":
		TrmId, err := jsonParsed.Path("trm_id").Data().(string)
		if !err {
			logger.Error("TerminalId field missing or not a string")
			return
		}
		h.Sessions.Add(conn, TrmId)
		logger.Info("Session opened:", TrmId)

	case "2":

		err := h.KafkaProducer.ProduceMessage(string(data))
		if err != nil {
			logger.Info("Failed to send message to Kafka", err)
		}
	default:
		logger.Warn("Unknown MsgType:", MsgType)
	}

	conn.Write([]byte("ACK\n"))
}

func (h *Handler) Close() {
	h.Sessions.CloseAll()
}
