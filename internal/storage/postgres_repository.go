package storage

import (
	"context"
	"sync"
	"time"
)

// TelemetryRecord represents persisted event data
type TelemetryRecord struct {
	EventID        string
	SourceClientID string
	EventType      string
	Timestamp      time.Time
	MetricValue    float64
	Attributes     map[string]string
}

// Repository abstracts database operations
type Repository interface {
	BatchInsert(ctx context.Context, records []TelemetryRecord) error
	GetLatestEvents(ctx context.Context, clientID string, limit int) ([]TelemetryRecord, error)
}

// MemoryPostgresRepository simulates PostgreSQL high-throughput bulk copy
type MemoryPostgresRepository struct {
	mu      sync.RWMutex
	records []TelemetryRecord
}

func NewMemoryPostgresRepository() *MemoryPostgresRepository {
	return &MemoryPostgresRepository{
		records: make([]TelemetryRecord, 0, 50000),
	}
}

func (r *MemoryPostgresRepository) BatchInsert(ctx context.Context, records []TelemetryRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, records...)
	return nil
}

func (r *MemoryPostgresRepository) GetLatestEvents(ctx context.Context, clientID string, limit int) ([]TelemetryRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []TelemetryRecord
	for i := len(r.records) - 1; i >= 0 && len(result) < limit; i-- {
		if r.records[i].SourceClientID == clientID {
			result = append(result, r.records[i])
		}
	}
	return result, nil
}
