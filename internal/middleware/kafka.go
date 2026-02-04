package middleware

import (
	"Payment-Terminal-Management/internal/config"
	"Payment-Terminal-Management/internal/utils"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func CreateKafkaProducer(kafkaCfg *config.KafkaProducerConfig) *kafka.Writer {
	log.Printf("Creating Kafka Producer with brokers: %v", kafkaCfg.BootstrapServers)
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaCfg.BootstrapServers...),
		Topic:                  kafkaCfg.ProducerTopic,
		RequiredAcks:           kafka.RequiredAcks(utils.ParseAcks(kafkaCfg.Acks)),
		BatchSize:              kafkaCfg.MaxInFlight,
		AllowAutoTopicCreation: true,
	}
}

func CreateKafkaConsumer(kafkaCfg *config.KafkaConsumerConfig) *kafka.Reader {
	log.Printf("Creating Kafka Consumer with brokers: %v, session timeout: %dms",
		kafkaCfg.BootstrapServers, kafkaCfg.HeartbeatInterval)

	// Convert milliseconds to duration
	sessionTimeout := time.Duration(60000) * time.Millisecond // 60s default
	heartbeatInterval := time.Duration(kafkaCfg.HeartbeatInterval) * time.Millisecond

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaCfg.BootstrapServers,
		GroupID: kafkaCfg.ConsumerGroupID,
		Topic:   kafkaCfg.ConsumerTopic,
		// Fetch behavior
		MinBytes: 1e3,
		MaxBytes: 10e6,
		MaxWait:  3 * time.Second, // Increased from 1s to 3s

		// Consumer group stability - use configured values
		SessionTimeout:    sessionTimeout,
		RebalanceTimeout:  90 * time.Second,
		HeartbeatInterval: heartbeatInterval,

		// Offset
		StartOffset: kafka.LastOffset, // Use latest offset

		// Commit
		CommitInterval: 0, // Manual commit

		// Isolation
		IsolationLevel: kafka.ReadCommitted,

		Logger:      nil,
		ErrorLogger: kafka.LoggerFunc(log.Printf),
	})
}
