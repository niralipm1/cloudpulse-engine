package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/niralipm1/cloudpulse-engine/internal/handlers"
	"github.com/niralipm1/cloudpulse-engine/internal/services"
	"github.com/niralipm1/cloudpulse-engine/internal/storage"
	"github.com/niralipm1/cloudpulse-engine/internal/stream"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println(" CloudPulse Distributed Telemetry Ingestion Engine")
	fmt.Println(" Author: Nirali Motka <niralplaystore1@gmail.com>")
	fmt.Println(" High-Throughput gRPC + Redis Streams Microservice")
	fmt.Println("==================================================")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	broker := stream.NewMemoryRedisBroker()
	defer broker.Close()

	repo := storage.NewMemoryPostgresRepository()
	processor := services.NewTelemetryProcessor(broker, repo)

	server := handlers.NewServer(9092, processor)
	if err := server.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server fatal error: %v\n", err)
		os.Exit(1)
	}
}
