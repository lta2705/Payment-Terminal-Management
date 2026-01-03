package handler

import (
	"Payment-Terminal-Management/internal/session"
	"Payment-Terminal-Management/internal/worker"
	"github.com/Jeffail/gabs"
	"github.com/bytedance/gopkg/util/logger"
)

type TransactionNotiHandler interface {
	Handle(raw []byte) error
}
type TransactionNotiHandlerIpml struct {
	session *session.SessionManager
	sender  worker.KafkaProducerWorker
}

func NewTransactionNotiHandler(session *session.SessionManager, sender worker.KafkaProducerWorker) TransactionNotiHandler {
	return &TransactionNotiHandlerIpml{
		session: session,
		sender:  sender,
	}
}

func (h *TransactionNotiHandlerIpml) Handle(raw []byte) error {
	jsonParsed, err := gabs.ParseJSON(raw)
	if err != nil {
		logger.Info("Failed to parse JSON:", err)
		return err
	}

	trmNode := jsonParsed.Path("TerminalId")
	if trmNode.Data() == nil {
		logger.Info("terminalId not found")
		return nil
	}

	trmId, ok := trmNode.Data().(string)
	if !ok || trmId == "" {
		logger.Info("terminalId invalid")
		return nil
	}

	if !h.session.Check(trmId) {
		logger.Info("Terminal not connected:", trmId)
		err := h.sender.SendMessage(string(raw))
		if err != nil {
			return err
		}
		return nil
	}

	ok = h.session.Send(trmId, string(raw))
	if !ok {
		logger.Info("Failed to send message to terminal:", trmId)
		sendErr := h.sender.SendMessage(string(raw))
		if sendErr != nil {
			return sendErr
		}
	}

	err2 := h.sender.SendMessage(string(raw))
	if err2 != nil {
		return err2
	}

	return nil
}
