package service

import (
	"Payment-Terminal-Management/internal/handler"
	"Payment-Terminal-Management/internal/session"
	"Payment-Terminal-Management/internal/worker"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type ConsumerService interface {
	StartConsumer(ctx context.Context)
}

type ConsumerServiceImpl struct {
	consumer worker.KafkaConsumerWorker
	session  *session.SessionManager
	handler  handler.TransactionNotiHandler
}

func (cs *ConsumerServiceImpl) StartConsumer(ctx context.Context) {
	logger.Info("Starting Consumer...")

	go cs.consumer.ConsumeMessage(func(msg kafka.Message) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		logger.Info("Message consumed", string(msg.Value))
		return cs.handler.Handle(msg.Value)
	})
}

func NewConsumerService(consumer worker.KafkaConsumerWorker, sessionManager *session.SessionManager, handler handler.TransactionNotiHandler) ConsumerService {
	return &ConsumerServiceImpl{
		consumer: consumer,
		session:  sessionManager,
		handler:  handler,
	}
}
