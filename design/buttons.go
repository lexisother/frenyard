package design

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)
//import "fmt"

type deUIDesignButton struct {
	framework.UILayoutProxy
	focusState float64
	overlay *framework.UIOverlayContainer
	button *framework.UIButton
	ripple deRippleFrame
	primary uint32
}

// FyFTick implements Frame.FyFTick
func (de *deUIDesignButton) FyFTick(delta float64) {
}

// FyETick overrides UILayoutProxy.FyETick
func (de *deUIDesignButton) FyETick(time float64) {
	de.UILayoutProxy.FyETick(time)
	if de.button.Hover || de.button.Down || de.button.Focused {
		de.focusState += time * 4
		if de.focusState >= 1 {
			de.focusState = 1
		}
	} else {
		de.focusState -= time * 4
		if de.focusState <= 0 {
			de.focusState = 0
		}
	}
	de.ripple.FyFTick(time)
}

// FyFDraw implements Frame.FyFDraw
func (de *deUIDesignButton) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass framework.FramePass) {
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
func (de *deUIDesignButton) FyFPadding() frenyard.Area2i {
	addedBorderX := sizeScale(16)
	// Don't completely ignore the subject but don't go doing anything silly either
	addedBorderY := sizeScale(4)
	return frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2})
}
// FyFClipping implements Frame.FyFClipping
func (de *deUIDesignButton) FyFClipping() bool {
	return true
}

// FyLSizeForLimits overrides UILayoutProxy.FyLSizeForLimits
func (de *deUIDesignButton) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	baseSize := de.UILayoutProxy.FyLSizeForLimits(limits)
	// The 36px is implemented as a min-height rather than a strict height to prevent incredible levels of what I can only refer to as "incredibly dumb, yet predictable, results".
	return baseSize.Max(frenyard.Vec2i{X: sizeScale(64), Y: sizeScale(36)})
}

func newDeUIDesignButtonPtr(primary uint32, content framework.UILayoutElement, behavior framework.ButtonBehavior) *framework.UIButton {
	des := &deUIDesignButton{
		primary: primary,
	}
	overlay := framework.NewUIOverlayContainerPtr(des, []framework.UILayoutElement{content})
	des.overlay = overlay
	framework.InitUILayoutProxy(des, overlay)
	des.button = framework.NewUIButtonPtr(des, behavior)
	des.ripple = deRippleFrame{
		Button: des.button,
		MaskRaw: borderButtonRaw,
		Scale: borderEffectiveScale,
	}
	return des.button
}
