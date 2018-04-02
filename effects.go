package animation

/*
Contains definitions of various animation effects that can be applied. These
effects should implement interface Animation
*/

import (
	"image/color"
	"time"
)

// Transition from one solid color (applied to all elements) to another solid
// color
type InterpolateSolid struct {
	startColor, endColor color.RGBA
	duration             time.Duration
	startTime            time.Time
}

func NewInterpolateSolid(startColor, endColor color.RGBA,
	duration time.Duration) InterpolateSolid {
	return InterpolateSolid{startColor, endColor, duration, time.Now()}
}

func (this *InterpolateSolid) Frame(buf []color.RGBA, frameTime time.Time) bool {
	return true
}
