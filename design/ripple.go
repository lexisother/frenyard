package design

import (
	"runtime"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
)

var rippleTex frenyard.Texture
func deSetupRipple() {
	rippleTex = integration.CreateHardcodedPNGTexture(rippleB64, []integration.ColourTransform{
		integration.ColourTransformBlueToStencil,
	})
}

// deRippleFrame implements the ripple effect as a frame.
// As ripples contain state, RippleFrame is stateful (cannot be shared between users), and requires ticking.
type deRippleFrame struct {
	Button *framework.UIButton
	MaskRaw framework.NinePatch
	Scale float64
	lastDown bool
	hasRipple bool
	downExpansionTimer float64
	hoverState float64
	ripplePosLock frenyard.Vec2i
}

// FyFTick implements Frame.FyFTick (though can also be used outside of it)
func (de *deRippleFrame) FyFTick(delta float64) {
	if de.Button.Hover || de.Button.Down || de.Button.Focused {
		de.hoverState += delta * 4
		if de.hoverState > 1 {
			de.hoverState = 1
		}
	} else {
		de.hoverState -= delta * 4
		if de.hoverState < 0 {
			de.hoverState = 0
		}
	}
	if de.hasRipple {
		if de.Button.Down {
			de.downExpansionTimer += delta
			if de.downExpansionTimer > 1 {
				de.downExpansionTimer = 1
			}
		} else {
			de.downExpansionTimer += delta * 2
			if de.downExpansionTimer >= 2 {
				de.hasRipple = false
			}
		}
	}
}

// FyFDraw implements Frame.FyFDraw (though can also be used outside of it)
func (de *deRippleFrame) FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass framework.FramePass) {
	if pass == framework.FramePassOverAfter {
		if de.Button.Down && !de.lastDown {
			de.downExpansionTimer = 0
			de.ripplePosLock = de.Button.LastMousePos
			de.hasRipple = true
		}
		de.lastDown = de.Button.Down
		areaSize := frenyard.Area2iOfSize(size)
		ripplePosition := de.ripplePosLock
		rippleSize := int32(de.downExpansionTimer * 2.5 * float64(frenyard.Max(size.X, size.Y)))

		drawColour := uint32(0x40FFFFFF)
		hoverColour := integration.ColourMix(0x00FFFFFF, 0x08FFFFFF, de.hoverState)
		if de.downExpansionTimer > 1 {
			drawColour = integration.ColourMix(drawColour, 0x00FFFFFF, de.downExpansionTimer - 1)
		}
		if de.hasRipple {
			tex := r.RenderToTexture(size, func () {
				r.Reset(0xFF000000)
				r.DrawRect(frenyard.DrawRectCommand{
					Tex: rippleTex,
					TexSprite: frenyard.Area2iOfSize(rippleTex.Size()),
					Target: frenyard.Area2iFromVecs(ripplePosition, frenyard.Vec2i{}).Align(frenyard.Vec2i{X: rippleSize, Y: rippleSize}, frenyard.Alignment2i{}),
					Colour: 0xFFFFFFFF,
				})
				de.MaskRaw.Draw(r, areaSize, de.Scale, frenyard.DrawRectCommand{
					Mode: frenyard.DrawModeModulate,
					Colour: 0xFFFFFFFF,
				})
			}, false)
			r.DrawRect(frenyard.DrawRectCommand{
				Tex: tex,
				TexSprite: areaSize,
				Target: areaSize,
				Colour: drawColour,
				Mode: frenyard.DrawModeAdd,
			})
			tex = nil
			runtime.GC()
		}
		de.MaskRaw.Draw(r, areaSize, de.Scale, frenyard.DrawRectCommand{
			Mode: frenyard.DrawModeAdd,
			Colour: hoverColour,
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
		MaskRaw: mask,
		Scale: scale,
	}
	assembledItem = framework.NewUIOverlayContainerPtr(rippleFrame, []framework.UILayoutElement{
		assembledItem,
	})
	rippleFrame.Button = framework.NewUIButtonPtr(assembledItem, click)
	return rippleFrame.Button
}
