package animation

// Code to support mapping from logical 'universes' to physical pixel layout.

import (
	"image/color"
)

// {board, strand, pixel} tuple identifying a physical pixel
type location struct {
	board, strand, pixel uint
}

// Structure to capture mapping from logical to physical layer
type Mapping struct {
	// Buffer of data mapping to physical pixels
	// Three levels of indexing:
	// 1. Controller board number
	// 2. Strand number within controller board
	// 3. Pixel number within strand
	physBuf [][][]color.RGBA

	// Mapping from 'universes' (logical view of pixels) to physical pixels.
	// Two levels of indexing:
	// 1. Universe number
	// 2. Pixel number within universe
	universes [][]location

	// Mapping from universe name to universe ID
	universeNameToId map[string]int
}

// Defines a range of physical pixels within asingle strand
type PhysicalRange struct {
	board, strand, startPixel, size uint
}

// Create a new Mapping, using the provided dimensions.
// Size of outer array governs the number of controller boards
// Sizes of inner arrays govern the number of strands within each board
// Values in inner array govern the number of pixels in the strand
func NewMapping(dimension [][]int) Mapping {
	// Make the triply-nested physical buffer structure based on the provided dimensions
	// Allocate space for a reasonable number of universes
	m := Mapping{make([][][]color.RGBA, len(dimension)), make([][]location, 10)[:0],
		make(map[string]int)}
	for boardIdx := range dimension {
		m.physBuf[boardIdx] = make([][]color.RGBA, len(dimension[boardIdx]))
		for strandIdx := range dimension[boardIdx] {
			m.physBuf[boardIdx][strandIdx] = make([]color.RGBA, dimension[boardIdx][strandIdx])
		}
	}
	return m
}

// Add a universe mapping with the given name.
// The provided set of physical ranges identifies the set of physical pixels
// corresponding to the universe. The order of physical pixels presented defines
// the logical ordering of the universe, and the size of the universe is equal
// to the number of physical pixels provided
// Returns true if the universe was successfully added; returns false if the
// universe name already exists or a specified physical pixel doesn't exist.
func (m *Mapping) AddUniverse(name string, ranges []PhysicalRange) bool {
	if _, exists := m.universeNameToId[name]; exists {
		return false
	}
	// Figure out the size
	size := uint(0)
	for _, r := range ranges {
		size += r.size
	}
	// Allocate locations array for universe
	locs := make([]location, size)
	// Populate locations array from pixel ranges
	unidx := 0
	for _, r := range ranges {
		for idx := r.startPixel; idx < r.startPixel+r.size; idx++ {
			locs[unidx] = location{r.board, r.strand, idx}
			unidx++
		}
	}

	// Add the universe to the structure
	m.universes = append(m.universes, locs)
	m.universeNameToId[name] = len(m.universes) - 1
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (m *Mapping) UpdateUniverse(universeId uint, universeData []color.RGBA) {
	u := m.universes[universeId]
	for idx, l := range u {
		m.physBuf[l.board][l.strand][l.pixel] = universeData[idx]
	}
}
