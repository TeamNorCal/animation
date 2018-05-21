package animation

import (
	"fmt"
	"image/color"
	"testing"
	"time"
)

func TestInterpolateSolid(t *testing.T) {
	effect := NewInterpolateSolid(color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}, time.Duration(1000))
	startTime := time.Now()
	effect.Start(startTime)

	buf := make([]color.RGBA, 1)
	result, done := effect.Frame(buf, startTime)

	val := result[0]
	if val.R != 0x00 || val.G != 0x00 || val.B != 0x00 || val.A != 0xff {
		t.Errorf("Value at start time not equal to start color: %v", val)
	}
	if done {
		t.Errorf("Done at start time")
	}

	result, done = effect.Frame(buf, startTime.Add(time.Duration(500)))
	fmt.Printf("Halfway point: %v\n", result[0])

	result, done = effect.Frame(buf, startTime.Add(time.Duration(1000)))
	val = result[0]
	if val.R != 0xff || val.G != 0xff || val.B != 0xff || val.A != 0xff {
		t.Errorf("Value at end time not equal to end color: %v", val)
	}
	if done {
		t.Errorf("Done on last real frame")
	}

	_, done = effect.Frame(buf, startTime.Add(time.Duration(10001)))
	if !done {
		t.Errorf("Not done after duration + 1")
	}
}
