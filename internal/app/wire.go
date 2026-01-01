//go:build wireinject
// +build wireinject

package app

import (
	"Payment-Terminal-Management/internal/config"
	"Payment-Terminal-Management/internal/middleware"
	"Payment-Terminal-Management/internal/service"
	"Payment-Terminal-Management/internal/session"
	"Payment-Terminal-Management/internal/transport"
	"Payment-Terminal-Management/internal/worker"
	"github.com/google/wire"
	"net"
	"os"
)

func ProvideListener() (net.Listener, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
	return net.Listen("tcp", ":"+port)
}

var ProducerSet = wire.NewSet(
	config.LoadKafkaProducerConfig,
	middleware.CreateKafkaProducer,
	worker.NewProducerWorker,
	service.NewProduceService,
	//wire.Bind(
	//	new(service.ProducerService),
	//	new(*service.ProducerServiceImpl),
	//),
)

var ConsumerSet = wire.NewSet(
	config.LoadKafkaConsumerConfig,
	middleware.CreateKafkaConsumer,
	worker.NewConsumerWorker,
	service.NewConsumerService,
)

var ServerSet = wire.NewSet(
	transport.NewServer,
	transport.NewHandler,
)

var sessionSet = wire.NewSet(
	session.NewSessionManager,
	wire.Bind(
		new(transport.SessionManager),
		new(*session.SessionManager),
	),
)

func InitializeApp() (*App, error) {
	wire.Build(
		ProvideListener,
		sessionSet,
		ConsumerSet,
		ProducerSet,
		ServerSet,
		NewApp,
	)
	return nil, nil
}
