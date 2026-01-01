package app

import (
	"Payment-Terminal-Management/internal/transport"
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

var sessionSet = wire.NewSet(
	NewSessionManager,
	wire.Bind(
		new(transport.SessionManager),
		new(*SessionManager),
	),
)

func InitializeApp() (*App, error) {
	wire.Build(
		ProvideListener,
		NewApp,
	)
	return &App{}, nil
}
