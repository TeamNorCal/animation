package animation

import (
	"image/color"
	"testing"
)

/*
Map of the universe defined below, with uinverse ID shown for each pixel:
(Board 0)
| 2 2 2 2 2 2 2 2 2 2
| . . . 3 . . . .
| . . . 1 1 1 1 . . . . . . . . . . . .

(Board 1)
| 2
| . . . . . . . . . . . . . . . . . . . . . . .
| . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . 3 3 3 3 1 1 1
| 1 1 1 1 1 . . . . . . . . . . . .

(Board 2)
| . . . . . . . . . . . . . . . . . . . .
| . 2 2 2 .
*/
func TestValidUniverse(t *testing.T) {
	// Set up some arbitrary strands across several board
	dimension := [][]int{
		[]int{10, 8, 19},
		[]int{1, 23, 64, 17},
		[]int{20, 5}}

	// Create a mapping
	mapping := NewMapping(dimension)

	if !mapping.AddUniverse("one", []PhysicalRange{
		PhysicalRange{Board: 0, Strand: 2, StartPixel: 3, Size: 4},
		PhysicalRange{Board: 1, Strand: 2, StartPixel: 61, Size: 3},
		PhysicalRange{Board: 1, Strand: 3, StartPixel: 0, Size: 5}}) {
		t.Error("Failed to add universe one")
	}

	if !mapping.AddUniverse("two", []PhysicalRange{
		PhysicalRange{Board: 1, Strand: 0, StartPixel: 0, Size: 1},
		PhysicalRange{Board: 0, Strand: 0, StartPixel: 0, Size: 10},
		PhysicalRange{Board: 2, Strand: 1, StartPixel: 1, Size: 3}}) {
		t.Error("Failed to add universe two")
	}

	if !mapping.AddUniverse("three", []PhysicalRange{
		PhysicalRange{Board: 0, Strand: 1, StartPixel: 3, Size: 1},
		PhysicalRange{Board: 1, Strand: 2, StartPixel: 57, Size: 4}}) {
		t.Error("Failed to add universe three")
	}

	id, err := mapping.IDForUniverse("one")
	if id != 0 || err != nil {
		t.Errorf("Unexpected universe ID %d for 'one' (error %v)", id, err)
	}

	id, err = mapping.IDForUniverse("two")
	if id != 1 || err != nil {
		t.Errorf("Unexpected universe ID %d for 'two' (error %v)", id, err)
	}

	id, err = mapping.IDForUniverse("three")
	if id != 2 || err != nil {
		t.Errorf("Unexpected universe ID %d for 'three' (error %v)", id, err)
	}

	oneData := make([]color.RGBA, 12)
	for idx := range oneData {
		oneData[idx] = color.RGBA{1, 1, 1, 1}
	}
	mapping.UpdateUniverse(0, oneData)

	twoData := make([]color.RGBA, 14)
	for idx := range twoData {
		twoData[idx] = color.RGBA{2, 2, 2, 2}
	}

	threeData := make([]color.RGBA, 5)
	for idx := range threeData {
		threeData[idx] = color.RGBA{3, 3, 3, 3}
	}

	mapping.UpdateUniverse(0, oneData)
	mapping.UpdateUniverse(1, twoData)
	mapping.UpdateUniverse(2, threeData)

	c0 := color.RGBA{0, 0, 0, 0}
	c1 := color.RGBA{1, 1, 1, 1}
	c2 := color.RGBA{2, 2, 2, 2}
	c3 := color.RGBA{3, 3, 3, 3}

	checkStrand(t, &mapping, 0, 0, func(idx int) color.RGBA {
		return c2
	})

	checkStrand(t, &mapping, 0, 1, func(idx int) color.RGBA {
		var expected color.RGBA
		switch {
		case idx < 3:
			expected = c0
		case idx == 3:
			expected = c3
		default:
			expected = c0
		}
		return expected
	})

	checkStrand(t, &mapping, 0, 2, func(idx int) color.RGBA {
		var expected color.RGBA
		switch {
		case idx < 3:
			expected = c0
		case idx < 7:
			expected = c1
		default:
			expected = c0
		}
		return expected
	})

	checkStrand(t, &mapping, 1, 0, func(idx int) color.RGBA {
		return c2
	})

	checkStrand(t, &mapping, 1, 1, func(idx int) color.RGBA {
		return c0
	})

	checkStrand(t, &mapping, 1, 2, func(idx int) color.RGBA {
		var expected color.RGBA
		switch {
		case idx < 57:
			expected = c0
		case idx < 61:
			expected = c3
		default:
			expected = c1
		}
		return expected
	})

	checkStrand(t, &mapping, 1, 3, func(idx int) color.RGBA {
		var expected color.RGBA
		switch {
		case idx < 5:
			expected = c1
		default:
			expected = c0
		}
		return expected
	})

	checkStrand(t, &mapping, 2, 0, func(idx int) color.RGBA {
		return c0
	})

	checkStrand(t, &mapping, 2, 1, func(idx int) color.RGBA {
		var expected color.RGBA
		switch {
		case idx < 1:
			expected = c0
		case idx < 4:
			expected = c2
		default:
			expected = c0
		}
		return expected
	})
}

func checkStrand(t *testing.T, mapping *Mapping, board, strand uint, expectedFunc func(idx int) color.RGBA) {
	data, err := mapping.GetStrandData(board, strand)
	if err != nil {
		t.Errorf("Error retrieving (%d, %d): %v", board, strand, err)
	}
	for idx, c := range data {
		expected := expectedFunc(idx)
		if expected != c {
			t.Errorf("Expected %v, got %v for strand (%d, %d) pixel %d", expected, c, board, strand, idx)
		}
	}
}
