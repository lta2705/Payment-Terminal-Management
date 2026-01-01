package main

import (
	"Payment-Terminal-Management/internal/app"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	application, err := app.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go application.Consumer.StartConsumer(ctx)

	// Run server
	go application.Start(ctx)

	<-ctx.Done()
	application.Stop()
	logger.Info("Application exited cleanly")
}
