package animation

/*
Contains definitions of various animation effects that can be applied. These
effects should implement interface Animation
*/

import (
	"image/color"
	"log"
	"os"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// InterpolateSolid transitions from one solid color (applied to all elements)
// to another solid color
type InterpolateSolid struct {
	startColor, endColor colorful.Color
	duration             time.Duration
	startTime            time.Time
	startOnCurrent       bool // Capture the current color and use it as the start color?
}

var fxlog = log.New(os.Stdout, "(EFFECT) ", 0)

func rgbaFromRGBHex(hexColor uint32) color.RGBA {
	return color.RGBA{uint8(hexColor >> 16 & 0xff), uint8(hexColor >> 8 & 0xff), uint8(hexColor & 0xff), 0xff}
}

// NewInterpolateSolidHexRGB creates an InterpolateSolid effect, given hex-encoded 24-bit RGB colors
func NewInterpolateSolidHexRGB(startColor, endColor uint32, duration time.Duration) *InterpolateSolid {
	startRGBA := rgbaFromRGBHex(startColor)
	endRGBA := rgbaFromRGBHex(endColor)
	return &InterpolateSolid{startColor: colorful.MakeColor(startRGBA), endColor: colorful.MakeColor(endRGBA), duration: duration}
}

// NewInterpolateSolid creates an InterpolateSolid effect
func NewInterpolateSolid(startColor, endColor color.RGBA,
	duration time.Duration) *InterpolateSolid {
	return &InterpolateSolid{startColor: colorful.MakeColor(startColor), endColor: colorful.MakeColor(endColor), duration: duration}
}

// NewInterpolateToHexRGB interpolates from the current color of the universe (determined by sampling the first element)
// to the provided end color, specified as a 24-bit RGB hex value
func NewInterpolateToHexRGB(endColor uint32, duration time.Duration) *InterpolateSolid {
	// Create a standard effect with arbitrary start color
	effect := NewInterpolateSolidHexRGB(0x0, endColor, duration)
	// ...then set the magic flag
	effect.startOnCurrent = true
	return effect
}

// Start starts the effect
func (effect *InterpolateSolid) Start(startTime time.Time) {
	fxlog.Printf("Setting start time %v", startTime)
	effect.startTime = startTime
}

// Frame generates an animation frame
func (effect *InterpolateSolid) Frame(buf []color.RGBA, frameTime time.Time) (output []color.RGBA, endSeq bool) {
	//fxlog.Printf("Buf cap: %d len: %d\n", cap(buf), len(buf))
	if frameTime.After(effect.startTime.Add(effect.duration)) {
		fxlog.Printf("Done at time %v (start time %v)\n", frameTime, effect.startTime)
		return buf, true
	}

	// See if we need to find the current universe color and use it as the start color
	if effect.startOnCurrent {
		sc := buf[0]
		sc.A = 0xff // Avoid a 0 transparency (in the case of an uninitialized buffer) which makes go-colorful unhappy
		effect.startColor = colorful.MakeColor(sc)
		effect.startOnCurrent = false // Clear the flag to prevent this from being done again
	}

	elapsed := frameTime.Sub(effect.startTime)
	completion := elapsed.Seconds() / effect.duration.Seconds()
	//fxlog.Printf("Frame at %2.2f%%", completion*100.0)
	//	currColorful := effect.startColor.BlendLab(effect.endColor, completion)
	currColorful := effect.startColor.BlendLuv(effect.endColor, completion)
	currColor := colorfulToRGBA(currColorful)
	for i := 0; i < len(buf); i++ {
		buf[i] = currColor
	}
	return buf, false
}

func colorfulToRGBA(c colorful.Color) color.RGBA {
	r, g, b := c.RGB255()
	return color.RGBA{r, g, b, 0xff}
}
