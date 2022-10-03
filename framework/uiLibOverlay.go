package framework

import (
	"github.com/uwu/frenyard"
)

/*
 * Details on this:
 * Frames are safe to share between containers if the frame itself is allowed to be shared between containers.
 * This becomes important for NinePatchFrames, introduced in ninepatch_frames.
 */

// FramePass is a draw pass for a Frame.
type FramePass uint8

// FramePassUnderBefore is Under, before the interior has been drawn.
const FramePassUnderBefore FramePass = 0

// FramePassUnderAfter is Under, after the interior has been drawn.
const FramePassUnderAfter FramePass = 1

// FramePassOverBefore is Over, before the interior has been drawn.
const FramePassOverBefore FramePass = 2

// FramePassOverAfter is Over, after the interior has been drawn.
const FramePassOverAfter FramePass = 3

// Frame represents a frame around the stacked elements in a UIOverlayContainer. It is, in a way, a simplified version of a single-element-container suitable for writing theme-related code.
type Frame interface {
	// FyFTick ticks the frame. Unstateful frames should really just ignore this.
	FyFTick(delta float64)
	// FyFDraw draws a pass of the frame.
	FyFDraw(r frenyard.Renderer, size frenyard.Vec2i, pass FramePass)
	// FyFPadding returns the padding. This is expected to not change without a notification - given by either resetting the holding container's frame or by doing something else that causes that to occur.
	FyFPadding() frenyard.Area2i
	// FyFClipping returns the requested PanelClipping value for when using this framing.
	FyFClipping() bool
}

// 'Overlay': Overlays elements over each other.
// Also an attachment point for a Frame.
// Useful for backgrounds.

// UIOverlayContainer overlays elements on top of each other, and this itself on top of a potentially padded NinePatchPackage.
type UIOverlayContainer struct {
	UIPanel
	UILayoutElementComponent
	_framing       Frame
	_state         []UILayoutElement
	_preferredSize frenyard.Vec2i
}

// NewUIOverlayContainerPtr creates a UIOverlayContainer
func NewUIOverlayContainerPtr(npp Frame, setup []UILayoutElement) *UIOverlayContainer {
	container := &UIOverlayContainer{
		UIPanel: NewPanel(frenyard.Vec2i{}),
	}
	InitUILayoutElementComponent(container)
	container.SetContent(npp, setup)
	container.FyEResize(container._preferredSize)
	return container
}

// FyLSubelementChanged implements UILayoutElement.FyLSubelementChanged
func (ufc *UIOverlayContainer) FyLSubelementChanged() {
	size := frenyard.Vec2i{}
	for _, v := range ufc._state {
		size = size.Max(v.FyLSizeForLimits(frenyard.Vec2iUnlimited()))
	}
	ufc._preferredSize = size.Add(ufc._framing.FyFPadding().Size())
	ufc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (ufc *UIOverlayContainer) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	if limits.Ge(ufc._preferredSize) {
		return ufc._preferredSize
	}
	max := frenyard.Vec2i{}
	paddingSize := ufc._framing.FyFPadding().Size()
	for _, v := range ufc._state {
		max = max.Max(v.FyLSizeForLimits(limits.Add(paddingSize.Negate())))
	}
	return max.Add(paddingSize)
}

// SetContent changes the content of the UIOverlayContainer.
func (ufc *UIOverlayContainer) SetContent(npp Frame, setup []UILayoutElement) {
	if ufc._state != nil {
		for _, v := range ufc._state {
			ufc.ThisUILayoutElementComponentDetails.Detach(v)
		}
	}
	ufc._state = setup
	ufc._framing = npp
	ufc.ThisUIPanelDetails.Clipping = npp.FyFClipping()
	for _, v := range setup {
		ufc.ThisUILayoutElementComponentDetails.Attach(v)
	}
	ufc.FyLSubelementChanged()
}

// FyEResize overrides UIPanel.FyEResize
func (ufc *UIOverlayContainer) FyEResize(size frenyard.Vec2i) {
	ufc.UIPanel.FyEResize(size)
	area := frenyard.Area2iOfSize(size).Contract(ufc._framing.FyFPadding())
	fixes := make([]PanelFixedElement, len(ufc._state))
	for idx, slot := range ufc._state {
		fixes[idx] = PanelFixedElement{
			Pos:     area.Pos(),
			Visible: true,
			Element: slot,
		}
		slot.FyEResize(area.Size())
	}

	ufc.ThisUIPanelDetails.SetContent(fixes)
}

// FyETick overrides UIPanel.FyETick
func (ufc *UIOverlayContainer) FyETick(delta float64) {
	ufc.UIPanel.FyETick(delta)
	ufc._framing.FyFTick(delta)
}

// FyEDraw overrides UIPanel.FyEDraw
func (ufc *UIOverlayContainer) FyEDraw(r frenyard.Renderer, under bool) {
	areaSize := ufc.FyESize()
	beforePass := FramePassOverBefore
	afterPass := FramePassOverAfter
	if under {
		beforePass = FramePassUnderBefore
		afterPass = FramePassUnderAfter
	}
	ufc._framing.FyFDraw(r, areaSize, beforePass)
	ufc.UIPanel.FyEDraw(r, under)
	ufc._framing.FyFDraw(r, areaSize, afterPass)
}
