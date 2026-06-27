package loadbalancer_test

import (
	"testing"

	"sentinel/internal/config"
	"sentinel/internal/loadbalancer"
)

func TestNew_Factory(t *testing.T) {
	tests := []struct {
		name     string
		strategy config.BalancerStrategy
		wantErr  bool
	}{
		{"default empty", "", false},
		{"round-robin", config.RoundRobin, false},
		{"random", config.Random, false},
		{"unknown", "invalid-strategy", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := loadbalancer.New(tt.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(%q) error = %v, wantErr %v", tt.strategy, err, tt.wantErr)
			}
			if !tt.wantErr && b == nil {
				t.Errorf("New(%q) returned nil balancer", tt.strategy)
			}
		})
	}
}
