package handlers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/niralipm1/cloudpulse-engine/internal/services"
)

// Server encapsulates the telemetry network listener and worker pool
type Server struct {
	processor *services.TelemetryProcessor
	port      int
}

func NewServer(port int, processor *services.TelemetryProcessor) *Server {
	return &Server{
		port:      port,
		processor: processor,
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	defer listener.Close()

	fmt.Printf("[CloudPulse Engine] gRPC Telemetry Server actively listening on %s\n", addr)

	// Simulate event ingestion lifecycle loop
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[CloudPulse Engine] Shutting down gracefully...")
			return nil
		case <-ticker.C:
			stats := s.processor.GetStats()
			fmt.Printf("[CloudPulse Telemetry] Processed: %d events | Errors: %d | p99: %.2fms\n",
				stats.TotalProcessed, stats.TotalErrors, stats.P99LatencyMs)
		}
	}
}
