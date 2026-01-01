package middleware

import (
	"Payment-Terminal-Management/internal/config"
	"Payment-Terminal-Management/internal/utils"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func CreateKafkaProducer(kafkaCfg *config.KafkaProducerConfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaCfg.BootstrapServers...),
		Topic:                  kafkaCfg.ProducerTopic,
		RequiredAcks:           kafka.RequiredAcks(utils.ParseAcks(kafkaCfg.Acks)),
		BatchSize:              kafkaCfg.MaxInFlight,
		AllowAutoTopicCreation: true,
	}
}

func CreateKafkaConsumer(kafkaCfg *config.KafkaConsumerConfig) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaCfg.BootstrapServers,
		GroupID: kafkaCfg.ConsumerGroupID,
		Topic:   kafkaCfg.ConsumerTopic,
		// Fetch behavior
		MinBytes: 1e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,

		// Consumer group stability
		SessionTimeout:   30 * time.Second,
		RebalanceTimeout: 60 * time.Second,

		// Offset
		StartOffset: kafka.FirstOffset,

		// Commit
		CommitInterval: 0,

		// Isolation
		IsolationLevel: kafka.ReadCommitted,

		Logger:      kafka.LoggerFunc(log.Printf),
		ErrorLogger: kafka.LoggerFunc(log.Printf),
	})
}
