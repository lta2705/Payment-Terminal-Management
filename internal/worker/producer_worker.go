package worker

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerWorker struct {
	Writer *kafka.Writer
}

func (pw *KafkaProducerWorker) SendMessage(key, value string) error {
	logger.Info("Message sent", key, value)
	return pw.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
}

func NewProducerWorker(writer *kafka.Writer) *KafkaProducerWorker {
	return &KafkaProducerWorker{
		Writer: writer,
	}
}
