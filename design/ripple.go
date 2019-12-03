package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
)

// deRippleFrame implements the ripple effect as a frame.
// As ripples contain state, RippleFrame is stateful (cannot be shared between users), and requires ticking.
type deRippleFrame struct {
	Button *framework.UIButton
	Mask framework.NinePatch
	Scale float64
}

// FyFTick implements Frame.FyFTick (though can also be used outside of it)
func (de *deRippleFrame) FyFTick(delta float64) {
	
}

// FyFDraw implements Frame.FyFDraw (though can also be used outside of it)
func (de *deRippleFrame) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass framework.FramePass) {
	if pass == framework.FramePassOverAfter {
		drawColour := uint32(0)
		if de.Button.Down {
			drawColour = 0x40FFFFFF
		} else if de.Button.Hover {
			drawColour = 0x0CFFFFFF
		} else if de.Button.Focused {
			drawColour = 0x08FFFFFF
		}
		de.Mask.Draw(r, frenyard.Area2iOfSize(size), de.Scale, frenyard.DrawRectCommand{
			Colour: drawColour,
		})
	}
}

// FyFPadding implements Frame.FyFPadding
func (de *deRippleFrame) FyFPadding() frenyard.Area2i {
	return frenyard.Area2i{}
}
// FyFClipping implements Frame.FyFClipping
func (de *deRippleFrame) FyFClipping() bool {
	return false
}

// deRippleInstall installs a ripple effect on a layout element.
func deRippleInstall(assembledItem framework.UILayoutElement, mask framework.NinePatch, scale float64, click framework.ButtonBehavior) framework.UILayoutElement {
	rippleFrame := &deRippleFrame{
		Mask: mask,
		Scale: scale,
	}
	assembledItem = framework.NewUIOverlayContainerPtr(rippleFrame, []framework.UILayoutElement{
		assembledItem,
	})
	rippleFrame.Button = framework.NewUIButtonPtr(assembledItem, click)
	return rippleFrame.Button
}
