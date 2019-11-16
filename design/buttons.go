package design

import "github.com/20kdc/CCUpdaterUI/frenyard"
//import "fmt"

type deUIDesignButton struct {
	frenyard.UILayoutProxy
	focusState float64
	attachedLabel *frenyard.UILabel
	overlay *frenyard.UIOverlayContainer
	button *frenyard.UIButton
}

// FyETick overrides UILayoutProxy.FyETick
func (de *deUIDesignButton) FyETick(time float64) {
	de.UILayoutProxy.FyETick(time)
	if de.button.Hover || de.button.Down {
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
	if de.attachedLabel != nil {
		if !de.button.Down {
			de.attachedLabel.SetColour(0xFF404040)
		} else {
			de.attachedLabel.SetColour(0xFFFFFFFF)
		}
	}
}

// FyFDraw implements Frame.FyFDraw
func (de *deUIDesignButton) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass frenyard.FramePass) {
	if pass == frenyard.FramePassUnderBefore {
		alpha := frenyard.ColourComponentClamp(int32(de.focusState * 255))
		alphaInv := 255 - alpha
		borderButtonShadow.Draw(r, frenyard.Area2iOfSize(size), borderEffectiveScale, frenyard.DrawRectCommand{
			Colour: (uint32(alphaInv) << 24) | 0xFFFFFF,
		})
		borderButtonShadowFocus.Draw(r, frenyard.Area2iOfSize(size), borderEffectiveScale, frenyard.DrawRectCommand{
			Colour: (uint32(alpha) << 24) | 0xFFFFFF,
		})
	} else if pass == frenyard.FramePassOverBefore {
		primaryColour := uint32(0xFF2040FF)
		if de.button.Down {
			primaryColour = 0xFF102080
		} else if de.button.Hover {
			primaryColour = 0xFF4060FF
		}
		borderButton.Draw(r, frenyard.Area2iOfSize(size), borderEffectiveScale, frenyard.DrawRectCommand{
			Colour: primaryColour,
		})
	}
}

// FyFPadding implements Frame.FyFPadding
func (de *deUIDesignButton) FyFPadding() frenyard.Area2i {
	addedBorderX := sizeScale(12)
	addedBorderY := sizeScale(8)
	return frenyard.Area2iFromVecs(frenyard.Vec2i{X: -addedBorderX, Y: -addedBorderY}, frenyard.Vec2i{X: addedBorderX * 2, Y: addedBorderY * 2})
}
// FyFClipping implements Frame.FyFClipping
func (de *deUIDesignButton) FyFClipping() bool {
	return true
}

func newDeUIDesignButtonPtr(content frenyard.UILayoutElement, label *frenyard.UILabel, behavior frenyard.ButtonBehavior) *frenyard.UIButton {
	des := &deUIDesignButton{}
	overlay := frenyard.NewUIOverlayContainerPtr(des, []frenyard.UILayoutElement{content})
	des.overlay = overlay
	frenyard.InitUILayoutProxy(des, overlay)
	des.button = frenyard.NewUIButtonPtr(des, behavior)
	return des.button
}
