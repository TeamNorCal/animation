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
