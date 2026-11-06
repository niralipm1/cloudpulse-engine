package services

import (
	"context"
	"testing"
	"time"

	"github.com/niralipm1/cloudpulse-engine/internal/storage"
	"github.com/niralipm1/cloudpulse-engine/internal/stream"
)

func TestTelemetryProcessor_ProcessEvent(t *testing.T) {
	broker := stream.NewMemoryRedisBroker()
	repo := storage.NewMemoryPostgresRepository()
	processor := NewTelemetryProcessor(broker, repo)
	ctx := context.Background()

	t.Run("valid event processed", func(t *testing.T) {
		input := TelemetryInput{
			EventID:        "evt-001",
			SourceClientID: "client-abc",
			EventType:      "heartbeat",
			Timestamp:      time.Now(),
			MetricValue:    42.0,
			Attributes:     map[string]string{"os": "android"},
		}

		msgID, err := processor.ProcessEvent(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if msgID == "" {
			t.Fatalf("expected non-empty msgID")
		}

		stats := processor.GetStats()
		if stats.TotalProcessed != 1 {
			t.Errorf("expected 1 processed event, got %d", stats.TotalProcessed)
		}
	})

	t.Run("missing client id returns error", func(t *testing.T) {
		input := TelemetryInput{
			EventID:   "evt-002",
			EventType: "heartbeat",
		}

		_, err := processor.ProcessEvent(ctx, input)
		if err == nil {
			t.Fatalf("expected error for missing client id")
		}
	})
}
