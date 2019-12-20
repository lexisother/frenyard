package framework

import (
	"fmt"
	"github.com/20kdc/CCUpdaterUI/frenyard"
)

/*
 * BEFORE CONTINUING: If you're here because of an "issue" where the scrollbox never activates the scrollbar,
 *  keep in mind that you need to let your layout crush the scrollbox for the scrollbox to actually do anything.
 * 
 * HTML notes for reference:
 * The Flexbox design is so utterly amazing that to prevent it from just outright ignoring size constraints in favour of minimum sizes, you have to actually just outright. Try it in your browser if you don't believe me:
 * <html><head></head><body style="display: flex;"><div>AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA</div><div>BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB</div><div>CAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAC</div></body></html>
 * This is with, of course, the default flex-shrink: 1 value.
 * The "fix" is "min-width: 0;" on the element that you want to fix,
 *  but the browser *refuses* to character-wrap unless given "overflow-wrap: break-word;" too.
 * Frenyard is somewhat less broken, and will by default implement this `min-width: 0;` technique.
 *  but labels actually character-wrap by default if they need to.
 * 
 * Anyways!
 * Judging by the official Material Design website, it seems like it leaves OS scrollboxes alone.
 * This is as good an excuse as any to use the "NinePatch this and use our own judgement" technology...
 */

// ScrollbarValue represents the value a scrollbar alters.
type ScrollbarValue interface {
	FySValue() float64
	// Sets the value. It is the responsibility of the target to clamp these values; this will show attempts at overscrolling for potential feedback.
	FySSetValue(v float64)
}

// ScrollbarTheme represents the Frame pair that makes up a Scrollbar.
type ScrollbarTheme struct {
	// The Base covers the whole visible main axis. The bounds of this touches the cross-axis sides of the Movement.
	Base Frame
	// The Movement covers the part of the scrollbar that moves during scrolling. It is given only the padding as bounds.
	Movement Frame
}

// UIScrollbar is an implementation of a scrollbar.
type UIScrollbar struct {
	UIElementComponent
	UILayoutElementComponent
	LayoutElementNoSubelementsComponent
	_theme ScrollbarTheme
	Value ScrollbarValue
	MouseNotch float64
	_vertical bool
	_grabbed bool
}

// NewUIScrollbarPtr creates a new scrollbar with the given theme and direction.
func NewUIScrollbarPtr(theme ScrollbarTheme, vertical bool, value ScrollbarValue) *UIScrollbar {
	uis := &UIScrollbar{
		_theme: theme,
		_vertical: vertical,
		Value: value,
		MouseNotch: 0.1,
	}
	InitUILayoutElementComponent(uis)
	return uis
}

// FyETick implements UIElement.FyETick
func (uis *UIScrollbar) FyETick(time float64) {
}

func (uis *UIScrollbar) _calcMetrics() (insideArea frenyard.Area2i, movementSize frenyard.Vec2i, movementArea frenyard.Area2i) {
	insideArea = frenyard.Area2iOfSize(uis.FyESize()).Contract(uis._theme.Base.FyFPadding())
	movementSize = uis._theme.Movement.FyFPadding().Size()
	movementArea = insideArea.Align(movementSize, frenyard.Alignment2i{})
	return insideArea, movementSize, movementArea
}

// FyEDraw implements UIElement.FyEDraw
func (uis *UIScrollbar) FyEDraw(r frenyard.Renderer, under bool) {
	if under {
		uis._theme.Base.FyFDraw(r, uis.FyESize(), FramePassUnderBefore)
	} else {
		uis._theme.Base.FyFDraw(r, uis.FyESize(), FramePassOverBefore)
	}
	// start movement
	insideArea, movementSize, movementArea := uis._calcMetrics()
	// Ok, here comes the "fun part"...
	if uis._vertical {
		movementArea.Y.Pos = insideArea.Y.Pos + int32(float64(insideArea.Y.Size - movementArea.Y.Size) * uis.Value.FySValue())
	} else {
		movementArea.X.Pos = insideArea.X.Pos + int32(float64(insideArea.X.Size - movementArea.X.Size) * uis.Value.FySValue())
	}
	r.Translate(movementArea.Pos())
	if under {
		uis._theme.Movement.FyFDraw(r, movementSize, FramePassUnderBefore)
		uis._theme.Movement.FyFDraw(r, movementSize, FramePassUnderAfter)
	} else {
		uis._theme.Movement.FyFDraw(r, movementSize, FramePassOverBefore)
		uis._theme.Movement.FyFDraw(r, movementSize, FramePassOverAfter)
	}
	r.Translate(movementArea.Pos().Negate())
	// end movement
	if under {
		uis._theme.Base.FyFDraw(r, uis.FyESize(), FramePassUnderAfter)
	} else {
		uis._theme.Base.FyFDraw(r, uis.FyESize(), FramePassOverAfter)
	}
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (uis *UIScrollbar) FyENormalEvent(me frenyard.NormalEvent) {
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (uis *UIScrollbar) FyEMouseEvent(me frenyard.MouseEvent) {
	decreaser, increaser := frenyard.MouseButtonScrollLeft, frenyard.MouseButtonScrollRight
	if uis._vertical {
		decreaser, increaser = frenyard.MouseButtonScrollUp, frenyard.MouseButtonScrollDown
	}
	if me.ID == frenyard.MouseEventDown && me.Button == increaser {
		uis.Value.FySSetValue(uis.Value.FySValue() + uis.MouseNotch)
	} else if me.ID == frenyard.MouseEventDown && me.Button == decreaser {
		uis.Value.FySSetValue(uis.Value.FySValue() - uis.MouseNotch)
	} else if me.ID == frenyard.MouseEventDown && me.Button == frenyard.MouseButtonLeft {
		uis._grabbed = true
	} else if me.ID == frenyard.MouseEventUp && me.Button == frenyard.MouseButtonLeft {
		uis._grabbed = false
	}
	if uis._grabbed {
		insideArea, _, _ := uis._calcMetrics()
		insideAreaMain := insideArea.X
		if uis._vertical {
			insideAreaMain = insideArea.Y
		}
		pointerMain := me.Pos.ConditionalTranspose(uis._vertical).X
		uis.Value.FySSetValue(float64(pointerMain - insideAreaMain.Pos) / float64(insideAreaMain.Size))
	}
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (uis *UIScrollbar) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	movementEffectiveSize := uis._theme.Movement.FyFPadding().Size()
	if uis._vertical {
		movementEffectiveSize.Y *= 2
	} else {
		movementEffectiveSize.X *= 2
	}
	return uis._theme.Base.FyFPadding().Size().Add(movementEffectiveSize)
}

// UIScrollbox allows scrolling around the contained element.
type UIScrollbox struct {
	UIPanel
	UILayoutElementComponent
	// The requirements for the scrollbar. Cached for simplicity's sake.
	_scrollbarMainCross frenyard.Vec2i
	_contained UILayoutElement
	_scrollbar *UIScrollbar
	// >0 means scrollbar is active
	_scrollLength int32
	_value float64
	_vertical bool
}

// NewUIScrollboxPtr creates a scrollbox.
func NewUIScrollboxPtr(theme ScrollbarTheme, content UILayoutElement, vertical bool) UILayoutElement {
	usc := &UIScrollbox{
		UIPanel: NewPanel(content.FyESize()),
		_contained: content,
		_vertical: vertical,
	}
	usc.ThisUIPanelDetails.Clipping = true
	usc._scrollbar = NewUIScrollbarPtr(theme, vertical, usc)
	InitUILayoutElementComponent(usc)
	usc.ThisUILayoutElementComponentDetails.Attach(usc._scrollbar)
	usc.ThisUILayoutElementComponentDetails.Attach(content)
	usc.FyEResize(usc.FyESize())
	return usc
}

// FyLSubelementChanged implements UILayoutElement.FyLSubelementChanged
func (usc *UIScrollbox) FyLSubelementChanged() {
	usc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (usc *UIScrollbox) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	limitsMainCross := limits.ConditionalTranspose(usc._vertical)
	baseSize := usc._contained.FyLSizeForLimits(limits)
	if baseSize.ConditionalTranspose(usc._vertical).X > limitsMainCross.X {
		// scrollbar penalty
		limitsNormalWithBar := frenyard.Vec2i{frenyard.SizeUnlimited, frenyard.AddCU(limitsMainCross.Y, -usc._scrollbarMainCross.X)}.ConditionalTranspose(usc._vertical)
		baseSize = usc._contained.FyLSizeForLimits(limitsNormalWithBar)
		baseSizeMainCross := baseSize.ConditionalTranspose(usc._vertical)
		return frenyard.Vec2i{usc._scrollbarMainCross.X, baseSizeMainCross.Y + usc._scrollbarMainCross.Y}.ConditionalTranspose(usc._vertical)
	}
	return baseSize
}

// FyEMouseEvent overrides UIPanel.FyEMouseEvent
func (usc *UIScrollbox) FyEMouseEvent(me frenyard.MouseEvent) {
	acceptance := false
	if usc._scrollLength > 0 && (me.ID == frenyard.MouseEventDown || me.ID == frenyard.MouseEventUp) {
		acceptance = me.Button == frenyard.MouseButtonScrollLeft || me.Button == frenyard.MouseButtonScrollRight
		if usc._vertical {
			acceptance = me.Button == frenyard.MouseButtonScrollDown || me.Button == frenyard.MouseButtonScrollUp
		}
	}
	if acceptance {
		usc._scrollbar.FyEMouseEvent(me)
	} else {
		usc.UIPanel.FyEMouseEvent(me)
	}
}

// FyEResize overrides UIPanel.FyEResize
func (usc *UIScrollbox) FyEResize(size frenyard.Vec2i) {
	usc.UIPanel.FyEResize(size)
	usc._scrollbarMainCross = usc._scrollbar.FyLSizeForLimits(frenyard.Vec2i{}).ConditionalTranspose(usc._vertical)

	baseSize := usc._contained.FyLSizeForLimits(size)
	baseSizeMainCross := baseSize.ConditionalTranspose(usc._vertical)
	sizeMainCross := size.ConditionalTranspose(usc._vertical)
	if baseSizeMainCross.X > sizeMainCross.X {
		if baseSizeMainCross.X == frenyard.SizeUnlimited {
			fmt.Printf("frenyard/UIScrollbox: Scroll layout was given element that wanted literally infinite size. This can't be made to work; track it down and get rid of it: %v %v\n", usc, baseSize)
		}
		// Subtract scrollbar from size; this is used for both layout and for calculating scroll len.
		sizeWithoutScrollbar := frenyard.Vec2i{baseSizeMainCross.X, sizeMainCross.Y - usc._scrollbarMainCross.Y}.ConditionalTranspose(usc._vertical)
		// Recalculate with the adjusted constraints because of the removed scrollbar
		baseSizeMainCross = usc._contained.FyLSizeForLimits(sizeWithoutScrollbar).ConditionalTranspose(usc._vertical)
		sizeWithoutScrollbar = frenyard.Vec2i{baseSizeMainCross.X, sizeMainCross.Y - usc._scrollbarMainCross.Y}.ConditionalTranspose(usc._vertical)
		usc._scrollLength = baseSizeMainCross.X - sizeMainCross.X
		usc._scrollbar.MouseNotch = (float64(sizeMainCross.X) / 6) / float64(usc._scrollLength)
		// can haz scrollbar?
		usc._scrollbar.FyEResize(frenyard.Vec2i{sizeMainCross.X, usc._scrollbarMainCross.Y}.ConditionalTranspose(usc._vertical))
		usc._contained.FyEResize(sizeWithoutScrollbar)
	} else {
		usc._scrollLength = 0
		usc._contained.FyEResize(size)
	}
	usc._updatePositions()
}

func (usc *UIScrollbox) _updatePositions() {
	if usc._scrollLength > 0 {
		scrollPos := -int32(usc._value * float64(usc._scrollLength))
		sizeMainCross := usc.FyESize().ConditionalTranspose(usc._vertical)
		usc.ThisUIPanelDetails.SetContent([]PanelFixedElement{
			PanelFixedElement{
				Element: usc._contained,
				Pos: frenyard.Vec2i{scrollPos, 0}.ConditionalTranspose(usc._vertical),
				Visible: true,
			},
			PanelFixedElement{
				Element: usc._scrollbar,
				Pos: frenyard.Vec2i{0, sizeMainCross.Y - usc._scrollbarMainCross.Y}.ConditionalTranspose(usc._vertical),
				Visible: true,
			},
		})
	} else {
		usc.ThisUIPanelDetails.SetContent([]PanelFixedElement{
			PanelFixedElement{
				Element: usc._contained,
				Visible: true,
			},
			// This needs to exist so that the structure doesn't change, for focus reasons (searchboxes)
			PanelFixedElement{
				Element: usc._scrollbar,
				Pos: frenyard.Vec2i{},
				Visible: false,
				Locked: true,
			},
		})
	}
}

// FySValue implements ScrollbarValue.FySValue
func (usc *UIScrollbox) FySValue() float64 {
	return usc._value
}

// FySSetValue implements ScrollbarValue.FySSetValue
func (usc *UIScrollbox) FySSetValue(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	usc._value = v
	usc._updatePositions()
}
