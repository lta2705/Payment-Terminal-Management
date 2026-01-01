package service

import (
	"Payment-Terminal-Management/internal/session"
	"Payment-Terminal-Management/internal/worker"
	"context"
	"github.com/Jeffail/gabs"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type ConsumerService interface {
	StartConsumer(ctx context.Context)
}

type ConsumerServiceImpl struct {
	consumer *worker.KafkaConsumerWorker
	session  *session.SessionManager
}

func (cs *ConsumerServiceImpl) StartConsumer(ctx context.Context) {
	logger.Info("Starting Consumer...")

	go cs.consumer.ConsumeMessage(func(msg kafka.Message) error {
		jsonParsed, err := gabs.ParseJSON(msg.Value)
		if err != nil {
			logger.Info("Failed to convert to struct:", err)
			return err
		}
		logger.Info("Message consumed", jsonParsed)

		trmNode := jsonParsed.Path("terminalId")
		if trmNode.Data() == nil {
			logger.Info("terminalId not found in message")
			return nil
		}

		trmId, ok := trmNode.Data().(string)
		if !ok || trmId == "" {
			logger.Info("terminalId invalid")
			return nil
		}
		if cs.session.Check(trmId) {
			cs.session.Send(trmId, string(msg.Value))
		}

		return nil
	})
}

func NewConsumerService(consumer *worker.KafkaConsumerWorker, sessionManager *session.SessionManager) ConsumerService {
	return &ConsumerServiceImpl{
		consumer: consumer,
		session:  sessionManager,
	}
}
