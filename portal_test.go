package animation

import (
	"testing"
	"time"
)

func TestPortal(t *testing.T) {
	p := NewPortal()
	resoStatus := make([]ResonatorStatus, 8)
	for idx := range resoStatus {
		resoStatus[idx] = ResonatorStatus{
			Health: 100.0,
			Level:  8,
		}
	}

	status := &PortalStatus{
		Faction:    ENL,
		Level:      8,
		Resonators: resoStatus,
	}

	p.UpdateStatus(status)
	p.GetFrame(time.Now())
}
