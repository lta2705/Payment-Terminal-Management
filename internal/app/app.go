package app

import (
	"Payment-Terminal-Management/internal/service"
	"Payment-Terminal-Management/internal/transport"
	"context"
	"net"
	"sync"
)

type App struct {
	Listener net.Listener
	Server   *transport.Server
	wg       sync.WaitGroup
	Consumer service.ConsumerService
}

func NewApp(listener net.Listener, server *transport.Server, Consumer service.ConsumerService) *App {
	return &App{
		Listener: listener,
		Server:   server,
		Consumer: Consumer,
	}
}

func (a *App) Start(ctx context.Context) {
	for {
		conn, err := a.Listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			a.Server.HandleConnection(ctx, conn)
		}()
	}
}

func (a *App) Stop() {
	a.Listener.Close()
	a.Server.Close()
	a.wg.Wait()
}
