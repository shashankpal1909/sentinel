package app

import (
	"time"

	"sentinel/internal/router"
)

type Snapshot struct {
	Runtime  *Runtime
	Router   *router.Router
	Version  uint64
	LoadedAt time.Time
}

// Current returns the snapshot itself, allowing Snapshot to satisfy SnapshotProvider interfaces directly.
func (s *Snapshot) Current() *Snapshot {
	return s
}
