package worker

import (
	"errors"
	"github.com/alitto/pond"
)

// PoolConfig holds the configuration for a worker pool.
type PoolConfig struct {
	Name             string // Name of the pool
	PrometheusPrefix string // Prefix for Prometheus metrics
	MinWorkers       uint16 // Minimum number of workers in the pool
	MaxWorkers       uint16 // Maximum number of workers in the pool
	ResizingStrategy string // Strategy for resizing the number of workers
	MaxQueuedJobs    uint16 // Maximum number of jobs that can be queued
}

// Default values for PoolConfig
const (
	DefaultName             = "default"
	DefaultPrometheusPrefix = "default"
	DefaultMinWorkers       = 4
	DefaultMaxWorkers       = 32
	DefaultResizingStrategy = "balanced"
	DefaultMaxQueuedJobs    = 100
)

// DefaultPoolConfig returns a PoolConfig with default values.
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		Name:             DefaultName,
		PrometheusPrefix: DefaultPrometheusPrefix,
		MinWorkers:       DefaultMinWorkers,
		MaxWorkers:       DefaultMaxWorkers,
		ResizingStrategy: DefaultResizingStrategy,
		MaxQueuedJobs:    DefaultMaxQueuedJobs,
	}
}

// resizerFromString returns a pond.ResizingStrategy based on the given name.
func resizerFromString(name string) (pond.ResizingStrategy, error) {
	switch name {
	case "eager":
		return pond.Eager(), nil
	case "lazy":
		return pond.Lazy(), nil
	case "balanced":
		return pond.Balanced(), nil
	default:
		return nil, errors.New("invalid resizer name")
	}
}
