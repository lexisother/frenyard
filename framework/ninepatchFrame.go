package framework

import "github.com/yellowsink/frenyard"

// A NinePatchFrameLayer provides the visuals for a given layer of a NinePatchFrame.
type NinePatchFrameLayer struct {
	Pass FramePass
	NinePatch
	Scale     float64
	ColourMod uint32
	Mode      frenyard.DrawMode
}

// A NinePatchFrame uses a set of NinePatches as a container.
type NinePatchFrame struct {
	Layers   []NinePatchFrameLayer
	Padding  frenyard.Area2i
	Clipping bool
}

// FyFDraw implements Frame.FyFDraw
func (npcf NinePatchFrame) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass FramePass) {
	for _, layer := range npcf.Layers {
		if layer.Pass != pass {
			continue
		}
		layer.Draw(r, frenyard.Area2iOfSize(size), layer.Scale, frenyard.DrawRectCommand{
			Colour: layer.ColourMod,
			Mode:   layer.Mode,
		})
	}
}

// FyFTick implements Frame.FyFTick
func (npcf NinePatchFrame) FyFTick(delta float64) {
}

// FyFPadding implements Frame.FyFPadding
func (npcf NinePatchFrame) FyFPadding() frenyard.Area2i {
	return npcf.Padding
}

// FyFClipping implements Frame.FyFClipping
func (npcf NinePatchFrame) FyFClipping() bool {
	return npcf.Clipping
}
