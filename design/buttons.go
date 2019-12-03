package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)
//import "fmt"

type deDesignButtonFrame struct {
	focusState float64
	button *framework.UIButton
	ripple deRippleFrame
	primary uint32
}

// FyFTick implements Frame.FyFTick
func (de *deDesignButtonFrame) FyFTick(delta float64) {
	if de.button.Hover || de.button.Down || de.button.Focused {
		de.focusState += delta * 4
		if de.focusState >= 1 {
			de.focusState = 1
		}
	} else {
		de.focusState -= delta * 4
		if de.focusState <= 0 {
			de.focusState = 0
		}
	}
	de.ripple.FyFTick(delta)
}

// FyFDraw implements Frame.FyFDraw
func (de *deDesignButtonFrame) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass framework.FramePass) {
	if pass == framework.FramePassUnderBefore {
		alpha := integration.ColourComponentClamp(int32(de.focusState * 255))
		alphaInv := 255 - alpha
		borderButtonShadow.Draw(r, frenyard.Area2iOfSize(size), borderEffectiveScale, frenyard.DrawRectCommand{
			Colour: (uint32(alphaInv) << 24) | 0xFFFFFF,
		})
		borderButtonShadowFocus.Draw(r, frenyard.Area2iOfSize(size), borderEffectiveScale, frenyard.DrawRectCommand{
			Colour: (uint32(alpha) << 24) | 0xFFFFFF,
		})
	} else if pass == framework.FramePassOverBefore {
		borderButton.Draw(r, frenyard.Area2iOfSize(size), borderEffectiveScale, frenyard.DrawRectCommand{
			Colour: de.primary,
		})
	}
	de.ripple.FyFDraw(r, size, pass)
}

// FyFPadding implements Frame.FyFPadding
func (de *deDesignButtonFrame) FyFPadding() frenyard.Area2i {
	addedBorderX := sizeScale(16)
	// Don't completely ignore the subject but don't go doing anything silly either
	addedBorderY := sizeScale(4)
	return frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2})
}
// FyFClipping implements Frame.FyFClipping
func (de *deDesignButtonFrame) FyFClipping() bool {
	return true
}

func newDeUIDesignButtonPtr(primary uint32, content framework.UILayoutElement, behavior framework.ButtonBehavior) *framework.UIButton {
	des := &deDesignButtonFrame{
		primary: primary,
	}
	minSizePanel := framework.NewPanel(frenyard.Vec2i{
		// The 36px is implemented as a min-height rather than a strict height to prevent incredible levels of what I can only refer to as "incredibly dumb, yet predictable, results".
		X: sizeScale(64),
		Y: sizeScale(36),
	})
	overlay := framework.NewUIOverlayContainerPtr(des, []framework.UILayoutElement{framework.ConvertElementToLayout(&minSizePanel), content})
	des.button = framework.NewUIButtonPtr(overlay, behavior)
	des.ripple = deRippleFrame{
		Button: des.button,
		MaskRaw: borderButtonRaw,
		Scale: borderEffectiveScale,
	}
	return des.button
}

type deCircleButtonFrame struct {
	focusState float64
	sizeState float64
	lastDown bool
	releasing bool
	button *framework.UIButton
}

// FyFTick implements Frame.FyFTick
func (de *deCircleButtonFrame) FyFTick(delta float64) {
	if de.releasing {
		de.sizeState += delta * 4
		de.focusState -= delta * 4
		if de.focusState <= 0 {
			de.focusState = 0
			de.sizeState = 1
			de.releasing = false
		}
	} else if de.button.Down {
		de.focusState += delta * 2
		if de.focusState >= 1 {
			de.focusState = 1
		}
	} else if de.lastDown && !de.button.Down {
		de.releasing = true
	} else if de.button.Hover || de.button.Focused {
		de.focusState += delta * 2
		if de.focusState >= 0.5 {
			de.focusState = 0.5
		}
	} else {
		de.focusState -= delta * 2
		if de.focusState <= 0 {
			de.focusState = 0
		}
	}
	if !de.releasing {
		de.sizeState = 1
	}
	de.lastDown = de.button.Down
}

// FyFDraw implements Frame.FyFDraw
func (de *deCircleButtonFrame) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass framework.FramePass) {
	if pass == framework.FramePassOverBefore {
		circleSize := int32(de.sizeState * float64(frenyard.Max(size.X, size.Y)))
		r.DrawRect(deEncodeCircleCmd(circleTexReductionsMask, frenyard.Vec2i{
			X: size.X / 2,
			Y: size.Y / 2,
		}, circleSize, frenyard.DrawRectCommand{
			Colour: integration.ColourMix(0x00FFFFFF, 0xFFFFFFFF, de.focusState * 0.5),
		}))
	}
}

// FyFPadding implements Frame.FyFPadding
func (de *deCircleButtonFrame) FyFPadding() frenyard.Area2i {
	return frenyard.Area2i{}
}
// FyFClipping implements Frame.FyFClipping
func (de *deCircleButtonFrame) FyFClipping() bool {
	return false
}

// Do be aware: Adds 8dp on each edge.
func newDeUICircleButtonPtr(content framework.UILayoutElement, behavior framework.ButtonBehavior) *framework.UIButton {
	frame := deCircleButtonFrame{
		focusState: 0.0,
	}
	content = framework.NewUIOverlayContainerPtr(&frame, []framework.UILayoutElement{content})
	button := framework.NewUIButtonPtr(content, behavior)
	frame.button = button
	return button
}
