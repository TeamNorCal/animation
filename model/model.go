package animationModel

import (
	"image/color"
)

// OpcChannel represents a channel in Open Pixel Controller parlance. Channel is a logical entity;
// the fcserver config file maps this to pixels on strands on particular FadeCandy boards
// fcserver configuration must honor this enumeration
type OpcChannel int

// ChannelData defines data for a particular OPC channel for a frame
type ChannelData struct {
	ChannelNum OpcChannel
	Data       []color.RGBA
}
