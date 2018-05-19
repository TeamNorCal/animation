package animation

import (
	"encoding/json"
	"image/color"
	"strconv"
	"time"
)

// Enacapsulates a model of a portal from the perpsective of animations.
// Provides an API semantically meaningful to Ingress

// Faction represents an Ingress fasction
type Faction int

const (
	// NEU means Neutral
	NEU Faction = iota
	// ENL means Enlightened!
	ENL
	// RES means Resistance :(
	RES
)

const windowSize = 30
const numResos = 8
const numShaftWindows = 16

// Universe defines a universe from the perspective of the animation engine.
// It consists of an index in the array of returned frame data, and a size
type Universe struct {
	Index int
	Size  int
}

// Universes defines universe data for the model. It's indexed by logical name,
// with values containing size of the universe and index into the array of
// frame data
var Universes map[string]Universe

// ResonatorStatus is the status of a resonator, you'll be surprised to hear
type ResonatorStatus struct {
	Level  int     // Resonator level, 0-8
	Health float32 // Resonator health, 0-100
}

// PortalStatus encapsulates the status of the portal
type PortalStatus struct {
	Faction    Faction           // Owning faction
	Level      float32           // Portal level, 0-8 (floating point, because average of resonator levels)
	Resonators []ResonatorStatus // Array of 8 resonators, level 0 (undeployed) - 8
}

// Portal encapsulates the animation status of the entire portal. This will probably be a singleton
// object, but the fields are encapsulated into a struct to allow for something different
type Portal struct {
	currentStatus *PortalStatus   // The cached current status of the portal
	sr            *SequenceRunner // SequenceRunner for portal portion
	resonators    []Animation     // Animations for resonators
	frameBuf      [][]color.RGBA  // Frame buffers by universe
}

// resonatorLevelColors is an array of colors of resonators of various levels, 0-8
var resonatorLevelColors = []uint32{
	0x000000, // L0
	0xEE8800, // L1
	0xFF6600, // L2
	0xCC3300, // L3
	0x990000, // L4
	0xFF0033, // L5
	0xCC0066, // L6
	0x660066, // L7
	0x330033, // L8
}

func init() {
	// Set up universes
	Universes = make(map[string]Universe)
	idx := 0
	// The resonators at the base
	for reso := 1; reso <= 8; reso++ {
		name := "base" + strconv.Itoa(reso)
		Universes[name] = Universe{
			Index: idx,
			Size:  windowSize,
		}
		idx++
	}

	// The, umm, shaft
	for level := 1; level <= 8; level++ {
		for window := 1; window <= 2; window++ {
			name := "towerLevel" + strconv.Itoa(level) + "Window" + strconv.Itoa(window)
			Universes[name] = Universe{
				Index: idx,
				Size:  windowSize,
			}
			idx++
		}
	}
}

// NewPortal creates a new portal structure
func NewPortal() *Portal {
	sizes := make([]uint, numShaftWindows)
	for idx := range sizes {
		sizes[idx] = windowSize
	}
	frameBuf := make([][]color.RGBA, numResos+numShaftWindows)
	for idx := range frameBuf {
		frameBuf[idx] = make([]color.RGBA, windowSize)
	}
	return &Portal{
		currentStatus: &PortalStatus{NEU, 0.0, make([]ResonatorStatus, numResos)},
		sr:            NewSequenceRunner(sizes),
		resonators:    make([]Animation, numResos),
		frameBuf:      frameBuf,
	}
}

// UpdateStatus updates the status of the portal from an animation perspective
func (p *Portal) UpdateStatus(status *PortalStatus) {
	newStatus := status.deepCopy()
	if p.currentStatus.Faction != newStatus.Faction || p.currentStatus.Level != newStatus.Level {
		p.updatePortal(status.Faction)
	}

	for idx, status := range newStatus.Resonators {
		if status != p.currentStatus.Resonators[idx] {
			p.updateResonator(idx, &status)
		}
	}

	p.currentStatus = newStatus
}

// GetFrame gets frame data for the portal, returning an array of frame data
// for each universe in the portal. Indices into this array are specified in the
// Universes map
// The returned buffers will typically be reused between frames, so callers
// should not hold onto references to them nor modify them!
func (p *Portal) GetFrame(frameTime time.Time) [][]color.RGBA {
	// Update resonators
	for idx := 0; idx < numResos; idx++ {
		p.getResoFrame(idx, frameTime)
	}
	p.sr.ProcessFrame(frameTime)
	for idx := 0; idx < numShaftWindows; idx++ {
		p.frameBuf[numResos+idx] = p.sr.UniverseData(uint(idx))
	}
	return p.frameBuf
}

func (msg *PortalStatus) deepCopy() (cpy *PortalStatus) {
	cpy = &PortalStatus{}

	byt, _ := json.Marshal(msg)
	json.Unmarshal(byt, cpy)
	return cpy
}

func (p *Portal) updatePortal(newFaction Faction) {

}

func (p *Portal) updateResonator(index int, newStatus *ResonatorStatus) {

}

// getResoFrame updates the frame buffer for the specified resonator with data
// for the current frame, with specified frame time
func (p *Portal) getResoFrame(index int, frameTime time.Time) {
	if p.resonators[index] == nil {
		return
	}
	buf, done := p.resonators[index].Frame(p.frameBuf[index], frameTime)
	p.frameBuf[index] = buf
	// Resonator animations run in a continuous loop, so restart if done
	if done {
		p.resonators[index].Start(frameTime)
	}
}
