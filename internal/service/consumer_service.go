package service

import (
	"github.com/lta2705/Payment-Terminal-Management/internal/worker"
)

type KafkaConsumerService struct {
	workers map[string]*KafkaConsumerWorker
}
