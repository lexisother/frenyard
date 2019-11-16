package frenyard

// A NinePatchFrameLayer provides the visuals for a given layer of a NinePatchFrame.
type NinePatchFrameLayer struct {
	NinePatch
	Scale float64
	ColourMod uint32
	Mode DrawMode
}

// A NinePatchFrame uses a set of NinePatches as a container.
type NinePatchFrame struct {
	UnderBefore NinePatchFrameLayer
	UnderAfter NinePatchFrameLayer
	OverBefore NinePatchFrameLayer
	OverAfter NinePatchFrameLayer
	Padding Area2i
	Clipping bool
}

// FyFDraw implements Frame.FyFDraw
func (npcf NinePatchFrame) FyFDraw(r Renderer, size Vec2i, pass FramePass) {
	layer := npcf.UnderBefore
	if pass == FramePassUnderAfter {
		layer = npcf.UnderAfter
	} else if pass == FramePassOverBefore {
		layer = npcf.OverBefore
	} else if pass == FramePassOverAfter {
		layer = npcf.OverAfter
	}
	layer.Draw(r, Area2iOfSize(size), layer.Scale, DrawRectCommand{
		Colour: layer.ColourMod,
		Mode: layer.Mode,
	})
}
// FyFPadding implements Frame.FyFPadding
func (npcf NinePatchFrame) FyFPadding() Area2i {
	return npcf.Padding
}
// FyFClipping implements Frame.FyFClipping
func (npcf NinePatchFrame) FyFClipping() bool {
	return npcf.Clipping
}
