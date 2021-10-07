package framework

import (
	"fmt"
	"github.com/yellowsink/frenyard"
)

// FocusEvent is an event type specific to the UI framework that represents focusing/unfocusing the receiving element.
type FocusEvent struct {
	// True if this was a focus, false if this was an unfocus.
	Focused bool
}

// FyVRoute implements NormalEvent.FyVRoute
func (fe FocusEvent) FyVRoute() frenyard.NormalEventRoute {
	return frenyard.NormalEventRouteStop
}

// FyVOffset implements NormalEvent.FyVOffset
func (fe FocusEvent) FyVOffset(o frenyard.Vec2i) frenyard.NormalEvent {
	return fe
}

/*
 * EXPLAINATION: EnterWindowEvent is sent when the root element is attached to the binding.
 * It's set to NormalEventRouteStructuralBroadcast. It is intended to go EVERYWHERE.
 * It is rebroadcast by panels when their contents change.
 * As panels and proxies are the only elements involved in event routing, this essentially means it is guaranteed to reach every element.
 */

// EnterWindowEvent is an event type specific to the UI framework that represents the element being attached to a window.
type EnterWindowEvent struct {
	Window frenyard.Window
}

// FyVRoute implements NormalEvent.FyVRoute
func (fe EnterWindowEvent) FyVRoute() frenyard.NormalEventRoute {
	return frenyard.NormalEventRouteStructuralBroadcast
}

// FyVOffset implements NormalEvent.FyVOffset
func (fe EnterWindowEvent) FyVOffset(o frenyard.Vec2i) frenyard.NormalEvent {
	return fe
}

/*
 * This is the core UIElement type without layout capabilities.
 * Simply put, if it's being drawn, it's this type.
 */

// UIElement is the core UI element type (no layout capabilities). An implementation must contain UIElementComponent or UIProxy.
type UIElement interface {
	FyENormalEvent(ev frenyard.NormalEvent)
	FyEMouseEvent(ev frenyard.MouseEvent)
	// Updates the element.
	FyETick(deltaTime float64)
	/*
	 * Drawing occurs in two passes: The 'under' pass and the main pass.
	 * If the element has a shadow, it draws that in the 'under' pass.
	 * As such, if an element has a background, it draws the shadow for that (if any) in the 'under' pass,
	 *  and splits its main pass into, in order: background, sub-element 'under' pass, sub-element main pass.
	 */
	FyEDraw(target frenyard.Renderer, under bool)
	/*
	 * Sets FyESize.
	 * FyESize MUST NOT change without FyEResize being used.
	 * FyEResize MUST ONLY be called if:
	 *  1. the parent element/window to which this element is attached is doing it
	 *  2. there is no parent element/window (setting default)
	 *  3. the parameter is equal to FyESize() (relayout)
	 * FyESize SHOULD default to a reasonable default size for the element.
	 */
	FyEResize(size frenyard.Vec2i)
	FyESize() frenyard.Vec2i
	// Attempts to find an element relative to Area2iOfSize(FyESize()). Does not check for itself. Returns an empty area for 'not found'.
	FyEFindElement(target UIElement) frenyard.Area2i
}

/*
 * A correct implementation of FyEResize & FyESize.
 * Part of core so it can't possibly get broken.
 */

// UIElementComponent implements the resizing logic for UIElement and default method implementations.
type UIElementComponent struct {
	// SUPER DUPER PRIVATE! DO NOT ACCESS OUTSIDE OF MEMBER METHODS.
	_fyUIElementComponentSize frenyard.Vec2i
}

// NewUIElementComponent creates a new UIElementComponent.
func NewUIElementComponent(size frenyard.Vec2i) UIElementComponent {
	return UIElementComponent{size}
}

// FyEResize implements UIElement.FyEResize
func (es *UIElementComponent) FyEResize(size frenyard.Vec2i) {
	es._fyUIElementComponentSize = size
}

// FyESize implements UIElement.FyESize
func (es *UIElementComponent) FyESize() frenyard.Vec2i {
	return es._fyUIElementComponentSize
}

// FyEFindElement implements UIElement.FyEFindElement
func (es *UIElementComponent) FyEFindElement(target UIElement) frenyard.Area2i {
	return frenyard.Area2i{}
}

type fyWindowElementBinding struct {
	window      frenyard.Window
	clearColour uint32
	element     UIElement
}

// CreateBoundWindow creates a window that is bound to an element.
func CreateBoundWindow(title string, vsync bool, clearColour uint32, e UIElement) (frenyard.Window, error) {
	return frenyard.GlobalBackend.CreateWindow(title, e.FyESize(), vsync, &fyWindowElementBinding{
		nil,
		clearColour,
		e,
	})
}

// FyRStart implements WindowReceiver.FyRStart
func (web *fyWindowElementBinding) FyRStart(w frenyard.Window) {
	web.window = w
	web.element.FyENormalEvent(EnterWindowEvent{w})
	web.element.FyENormalEvent(FocusEvent{true})
}

// FyRTick implements WindowReceiver.FyRTick
func (web *fyWindowElementBinding) FyRTick(f float64) {
	if !web.window.Size().Eq(web.element.FyESize()) {
		web.element.FyEResize(web.window.Size())
	}
	web.element.FyETick(f)
	web.window.Reset(web.clearColour)
	web.element.FyEDraw(web.window, true)
	web.element.FyEDraw(web.window, false)
	web.window.Present()
}

// FyRNormalEvent implements WindowReceiver.FyRNormalEvent
func (web *fyWindowElementBinding) FyRNormalEvent(ev frenyard.NormalEvent) {
	web.element.FyENormalEvent(ev)
}

// FyRMouseEvent implements WindowReceiver.FyRMouseEvent
func (web *fyWindowElementBinding) FyRMouseEvent(ev frenyard.MouseEvent) {
	web.element.FyEMouseEvent(ev)
}

// FyRClose implements WindowReceiver.FyRClose
func (web *fyWindowElementBinding) FyRClose() {
}

// PanelFixedElement describes an element attached to a panel.
type PanelFixedElement struct {
	Pos frenyard.Vec2i
	// Setting this to false is useful if you want an element to still tick but want to remove the drawing overhead.
	Visible bool
	// Setting this to true 'locks' the element. The element still participates in hit-tests but fails to focus and events are NOT forwarded.
	Locked  bool
	Element UIElement
}

/*
 * Basic "set it and forget it" stateful panel that does not transmit or receive layout data.
 * This is part of core because it's responsible for implementing several UI rules, such as cursor forwarding.
 */

// UIPanel is a "set it and forget it" stateful panel for placing multiple elements into.
type UIPanel struct {
	UIElementComponent
	ThisUIPanelDetails UIPanelDetails
}

// UIPanelDetails contains the details of a UIPanel otherwise accessible only by it's owner.
type UIPanelDetails struct {
	// Enables/disables clipping
	Clipping bool
	// This is a bitfield
	_buttonsDown uint16
	// Mouse event receiver. Be aware: Focus outside of Content is not a very good idea, except -1 (None)
	_focus int
	// Content (As far as I can tell there is no way to change the length of a slice without replacing it.)
	_content []PanelFixedElement
	// Cached window
	_window frenyard.Window
}

// NewPanel creates a UIPanel.
func NewPanel(size frenyard.Vec2i) UIPanel {
	return UIPanel{
		NewUIElementComponent(size),
		UIPanelDetails{
			false,
			0,
			-1,
			make([]PanelFixedElement, 0),
			nil,
		},
	}
}

// SetContent sets the contents of the panel.
func (pan *UIPanelDetails) SetContent(content []PanelFixedElement) {
	// Before we touch anything, broadcast EnterWindowEvent. Panels that were already in the same window will drop this.
	if pan._window != nil {
		for _, v := range content {
			v.Element.FyENormalEvent(EnterWindowEvent{
				pan._window,
			})
		}
	}

	// Is this actually a change we need to worry about?
	// DO BE WARNED: THIS IS A LOAD-BEARING OPTIMIZATION. DISABLE IT AND BUTTONS DON'T WORK PROPERLY
	// Reason: Clicking a button changes the button content which causes a layout rebuild.
	// Layout rebuilds destroying focus also destroys the evidence the button was pressed.
	changeCanBeIgnored := true
	if len(content) != len(pan._content) {
		changeCanBeIgnored = false
	} else {
		// Lengths are the same; if the elements are the same, we can just roll with it
		for k, v := range content {
			if pan._content[k].Element != v.Element {
				changeCanBeIgnored = false
			}
		}
	}
	if !changeCanBeIgnored {
		if pan._focus != -1 {
			// Ensure the focus has been notified.
			focusElement := pan._content[pan._focus]
			// Has to occur before the buttons get removed or ordering issues occur.
			pan._focus = -1
			// And we've successfully delivered the MOUSEDOWNs to the *new* element, -1, by default
			for button := (uint)(0); button < (uint)(frenyard.MouseButtonLength); button++ {
				if pan._buttonsDown&(1<<button) != 0 {
					focusElement.Element.FyEMouseEvent(frenyard.MouseEvent{
						frenyard.Vec2i{0, 0},
						frenyard.MouseEventUp,
						(frenyard.MouseButton)(button),
					})
				}
			}
			focusElement.Element.FyENormalEvent(FocusEvent{false})
		}
	}
	pan._content = content
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (pan *UIPanel) FyENormalEvent(ev frenyard.NormalEvent) {
	switch xev := ev.(type) {
	case EnterWindowEvent:
		if pan.ThisUIPanelDetails._window == xev.Window {
			// Drop the event, it's redundant
			return
		}
		pan.ThisUIPanelDetails._window = xev.Window
	}
	route := ev.FyVRoute()
	if (route == frenyard.NormalEventRouteBroadcast) || (route == frenyard.NormalEventRouteStructuralBroadcast) {
		for _, v := range pan.ThisUIPanelDetails._content {
			if (route == frenyard.NormalEventRouteStructuralBroadcast) || (v.Visible && !v.Locked) {
				v.Element.FyENormalEvent(ev.FyVOffset(v.Pos.Negate()))
			}
		}
	} else {
		if pan.ThisUIPanelDetails._focus != -1 {
			elem := pan.ThisUIPanelDetails._content[pan.ThisUIPanelDetails._focus]
			if elem.Visible && !elem.Locked {
				elem.Element.FyENormalEvent(ev)
			}
		}
	}
}

func (pan *UIPanel) _fyUIPanelForwardMouseEvent(target PanelFixedElement, ev frenyard.MouseEvent) {
	ev = ev.Offset(target.Pos.Negate())
	// Problematic mouse events are prevented from reaching locked targets via the hit-test logic.
	target.Element.FyEMouseEvent(ev)
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (pan *UIPanel) FyEMouseEvent(ev frenyard.MouseEvent) {
	// Useful for debugging if any of the warnings come up
	// if ev.ID != MouseEventMove { fmt.Printf("ui_core.go/Panel (%p)/FyEMouseEvent %v %v (%v, %v)\n", pan, ev.ID, ev.Button, ev.Pos.X, ev.Pos.Y) }

	invalid := false
	hittest := -1
	buttonMask := (uint16)(0)
	if ev.Button != -1 {
		buttonMask = (uint16)(1 << (uint)(ev.Button))
	}
	// Hit-test goes in reverse so that the element drawn last wins.
	for keyRev := range pan.ThisUIPanelDetails._content {
		key := len(pan.ThisUIPanelDetails._content) - (keyRev + 1)
		val := pan.ThisUIPanelDetails._content[key]
		if !val.Visible {
			continue
		}
		if frenyard.Area2iFromVecs(val.Pos, val.Element.FyESize()).Contains(ev.Pos) {
			//fmt.Printf(" Hit index %v\n", key)
			hittest = key
			if val.Locked {
				hittest = -1
			}
			break
		}
	}
	switch ev.ID {
	case frenyard.MouseEventMove:
		// Mouse-move events go everywhere.
		for _, val := range pan.ThisUIPanelDetails._content {
			pan._fyUIPanelForwardMouseEvent(val, ev)
		}
		invalid = true
	case frenyard.MouseEventUp:
		if pan.ThisUIPanelDetails._buttonsDown&buttonMask == 0 {
			fmt.Printf("ui_core.go/Panel (%p)/FyEMouseEvent warning: Button removal on non-existent button %v\n", pan, ev.Button)
			invalid = true
		} else {
			pan.ThisUIPanelDetails._buttonsDown &= 0xFFFF ^ buttonMask
		}
	case frenyard.MouseEventDown:
		if pan.ThisUIPanelDetails._buttonsDown == 0 {
			/*
			 * FOCUS REASONING DESCRIPTION
			 * If focusing on a subelement of an unfocused panel
			 *  the parent focuses the panel
			 *  the panel gets & forwards focus message to old interior focus
			 *  the panel gets the mouse event
			 *  the panel creates & forwards unfocus message to old interior focus
			 *  the panel creates & forwards focus message to new interior focus
			 * If changing the subelement of a focused panel
			 *  the panel creates & forwards unfocus message to old interior focus
			 *  the panel creates & forwards focus message to new interior focus
			 * If unfocusing a panel
			 *  the panel gets & forwards unfocus message to interior focus
			 */
			if pan.ThisUIPanelDetails._focus != hittest {
				// Note that this only happens when all other buttons have been released.
				// This prevents having to create fake release events.
				// The details of the order here are to do with issues when elements start modifying things in reaction to events.
				// Hence, the element that is being focused gets to run first so it will always receive an unfocus event after it has been focused.
				// While the element being unfocused is unlikely to get refocused under sane circumstances.
				// If worst comes to worst, make this stop sending focus events so nobody has to worry about focus state atomicity.
				oldFocus := pan.ThisUIPanelDetails._focus
				pan.ThisUIPanelDetails._focus = hittest
				newFocusFixed := PanelFixedElement{}
				if pan.ThisUIPanelDetails._focus != -1 {
					newFocusFixed = pan.ThisUIPanelDetails._content[pan.ThisUIPanelDetails._focus]
				}
				// Since a mouse event came in in the first place, we know the panel's focused.
				// Focus the newly focused element.
				if newFocusFixed.Element != nil {
					newFocusFixed.Element.FyENormalEvent(FocusEvent{true})
				}
				// Unfocus the existing focused element, if any.
				if oldFocus != -1 {
					pan.ThisUIPanelDetails._content[oldFocus].Element.FyENormalEvent(FocusEvent{false})
				}
			}
		}
		if pan.ThisUIPanelDetails._buttonsDown&buttonMask != 0 {
			fmt.Println("ui_core.go/Panel/FyEMouseEvent warning: Button added when it was already added")
			invalid = true
		} else {
			pan.ThisUIPanelDetails._buttonsDown |= buttonMask
		}
	}
	// Yes, focus gets to receive mouse-move events out of bounds even if there are no buttons.
	// All the state is updated, forward the event
	if !invalid && pan.ThisUIPanelDetails._focus != -1 {
		pan._fyUIPanelForwardMouseEvent(pan.ThisUIPanelDetails._content[pan.ThisUIPanelDetails._focus], ev)
	}
}

// FyEDraw implements UIElement.FyEDraw
func (pan *UIPanel) FyEDraw(target frenyard.Renderer, under bool) {
	if pan.ThisUIPanelDetails.Clipping {
		// Clipping: everything is inside panel bounds
		if under {
			return
		}
		oldClip := target.Clip()
		newClip := oldClip.Intersect(frenyard.Area2iOfSize(pan.FyESize()))
		if newClip.Empty() {
			return
		}
		target.SetClip(newClip)
		defer target.SetClip(oldClip)
		for pass := 0; pass < 2; pass++ {
			for _, val := range pan.ThisUIPanelDetails._content {
				if !val.Visible {
					continue
				}
				target.Translate(val.Pos)
				val.Element.FyEDraw(target, pass == 0)
				target.Translate(val.Pos.Negate())
			}
		}
	} else {
		// Not clipping; this simply arranges a bunch of elements
		for _, val := range pan.ThisUIPanelDetails._content {
			if !val.Visible {
				continue
			}
			target.Translate(val.Pos)
			val.Element.FyEDraw(target, under)
			target.Translate(val.Pos.Negate())
		}
	}
}

// FyETick implements UIElement.FyETick
func (pan *UIPanel) FyETick(f float64) {
	for _, val := range pan.ThisUIPanelDetails._content {
		if !val.Visible {
			continue
		}
		val.Element.FyETick(f)
	}
}

// FyEFindElement implements UIElement.FyEFindElement
func (pan *UIPanel) FyEFindElement(target UIElement) frenyard.Area2i {
	for _, val := range pan.ThisUIPanelDetails._content {
		if val.Element == target {
			return frenyard.Area2iFromVecs(val.Pos, val.Element.FyESize())
		}
		area := val.Element.FyEFindElement(target)
		if !area.Empty() {
			return area.Translate(val.Pos)
		}
	}
	return frenyard.Area2i{}
}

// UIProxyHost is used to 'drill down' to the UIProxy within an element.
type UIProxyHost interface {
	// Returns the *UIProxy within this element.
	fyGetUIProxy() *UIProxy
}

// UIProxy is a "proxy" element. Useful to use another element as a base class without including it via inheritance.
type UIProxy struct {
	// This element is semi-private: it may be read by UIProxy and UILayoutProxy but nothing else.
	fyUIProxyTarget UIElement
}

func (px *UIProxy) fyGetUIProxy() *UIProxy {
	return px
}

// InitUIProxy initializes a UIProxy, setting the target.
func InitUIProxy(proxy UIProxyHost, target UIElement) {
	proxy.fyGetUIProxy().fyUIProxyTarget = target
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (px *UIProxy) FyENormalEvent(ev frenyard.NormalEvent) {
	px.fyUIProxyTarget.FyENormalEvent(ev)
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (px *UIProxy) FyEMouseEvent(ev frenyard.MouseEvent) {
	px.fyUIProxyTarget.FyEMouseEvent(ev)
}

// FyEDraw implements UIElement.FyEDraw
func (px *UIProxy) FyEDraw(target frenyard.Renderer, under bool) {
	px.fyUIProxyTarget.FyEDraw(target, under)
}

// FyETick implements UIElement.FyETick
func (px *UIProxy) FyETick(f float64) {
	px.fyUIProxyTarget.FyETick(f)
}

// FyEResize implements UIElement.FyEResize
func (px *UIProxy) FyEResize(v frenyard.Vec2i) {
	px.fyUIProxyTarget.FyEResize(v)
}

// FyESize implements UIElement.FyESize
func (px *UIProxy) FyESize() frenyard.Vec2i {
	return px.fyUIProxyTarget.FyESize()
}

// FyEFindElement implements UIElement.FyEFindElement
func (px *UIProxy) FyEFindElement(target UIElement) frenyard.Area2i {
	if px.fyUIProxyTarget == target {
		return frenyard.Area2iOfSize(px.FyESize())
	}
	return px.fyUIProxyTarget.FyEFindElement(target)
}
