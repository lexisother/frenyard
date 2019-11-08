/*
 * PLEASE KEEP IN MIND: NONE OF THIS IS 'PUBLIC' FOR DEPENDENCY PURPOSES (yet)
 */

package frenyard

// For warnings
import "fmt"

type FocusEvent struct {
	Focused bool
}

/*
 * This is the core UIElement type without layout capabilities.
 * Simply put, if it's being drawn, it's this type.
 */
type UIElement interface {
	FyENormalEvent(ev NormalEvent)
	FyEMouseEvent(ev MouseEvent)
	// Updates the element.
	FyETick(deltaTime float64)
	/*
	 * Drawing occurs in two passes: The 'under' pass and the main pass.
	 * If the element has a shadow, it draws that in the 'under' pass.
	 * As such, if an element has a background, it draws the shadow for that (if any) in the 'under' pass,
	 *  and splits its main pass into, in order: background, sub-element 'under' pass, sub-element main pass.
	 */
	FyEDraw(target Renderer, under bool)
	/*
	 * Sets FyESize.
	 * FyESize MUST NOT change without FyEResize being used.
	 * FyEResize MUST ONLY be called if:
	 *  1. the parent element/window to which this element is attached is doing it
	 *  2. there is no parent element/window (setting default)
	 *  3. the parameter is equal to FyESize() (relayout)
	 * FyESize SHOULD default to a reasonable default size for the element.
	 */
	FyEResize(size Vec2i)
	FyESize() Vec2i
}

/*
 * A correct implementation of FyEResize & FyESize.
 * Part of core so it can't possibly get broken.
 */
type UIElementComponent struct {
	// SUPER DUPER PRIVATE! DO NOT ACCESS OUTSIDE OF MEMBER METHODS.
	_fy_UIElementComponent_size Vec2i
}
func NewUIElementComponent(size Vec2i) UIElementComponent {
	return UIElementComponent{size}
}
func (es *UIElementComponent) FyEResize(size Vec2i) {
	es._fy_UIElementComponent_size = size
}
func (es *UIElementComponent) FyESize() Vec2i {
	return es._fy_UIElementComponent_size
}

type fyWindowElementBinding struct {
	window Window
	clearColour uint32
	element UIElement
}
func CreateBoundWindow(title string, vsync bool, clearColour uint32, e UIElement) (Window, error) {
	return GlobalBackend.CreateWindow(title, e.FyESize(), vsync, &fyWindowElementBinding{
		nil,
		clearColour,
		e,
	})
}
func (web *fyWindowElementBinding) FyRStart(w Window) {
	web.window = w
	web.element.FyENormalEvent(FocusEvent{true})
}
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
func (web *fyWindowElementBinding) FyRNormalEvent(ev NormalEvent) {
	web.element.FyENormalEvent(ev)
}
func (web *fyWindowElementBinding) FyRMouseEvent(ev MouseEvent) {
	web.element.FyEMouseEvent(ev)
}
func (web *fyWindowElementBinding) FyRClose() {
	web.window.Destroy()
}

/*
 * Type used by Panel for attached elements.
 */
type PanelFixedElement struct {
	Pos Vec2i
	// Setting this to false is useful if you want an element to still tick but want to remove the drawing overhead.
	Visible bool
	Element UIElement
}
func NewPanelFixedElement(pos Vec2i, elem UIElement) PanelFixedElement {
	return PanelFixedElement{
		pos,
		true,
		elem,
	}
}

/*
 * Basic "set it and forget it" stateful panel that does not transmit or receive layout data.
 * This is part of core because it's responsible for implementing several UI rules.
 */
type UIPanel struct {
	UIElementComponent
	// Enables/disables clipping
	PanelClipping bool
	_fy_UIPanel fyUIPanelPrivate
}
type fyUIPanelPrivate struct {
	// This is a bitfield
	buttonsDown uint16
	// Mouse event receiver. Be aware: Focus outside of Content is not a very good idea, except -1 (None)
	focus int
	// Content (As far as I can tell there is no way to change the length of a slice without replacing it.)
	content []PanelFixedElement
}

func NewPanel(size Vec2i) UIPanel {
	return UIPanel{
		NewUIElementComponent(size),
		false,
		fyUIPanelPrivate{
			0,
			-1,
			make([]PanelFixedElement, 0),
		},
	}
}

func (pan *UIPanel) UIPanelSetContent(content []PanelFixedElement) {
	// Is this actually a change we need to worry about?
	// DO BE WARNED: THIS IS A LOAD-BEARING OPTIMIZATION. DISABLE IT AND BUTTONS DON'T WORK PROPERLY
	// Reason: Clicking a button changes the button content which causes a layout rebuild.
	// Layout rebuilds destroying focus also destroys the evidence the button was pressed.
	changeCanBeIgnored := true
	if len(content) != len(pan._fy_UIPanel.content) {
		changeCanBeIgnored = false
	} else {
		// Lengths are the same; if the elements are the same, we can just roll with it
		for k, v := range content {
			if pan._fy_UIPanel.content[k].Element != v.Element {
				changeCanBeIgnored = false
			}
		}
	}
	if !changeCanBeIgnored {
		if pan._fy_UIPanel.focus != -1 {
			// Ensure the focus has been notified.
			focusElement := pan._fy_UIPanel.content[pan._fy_UIPanel.focus]
			// And we've successfully delivered the MOUSEDOWNs to the *new* element, -1, by default
			for button := (int8)(0); button < MOUSEBUTTON_LENGTH; button++ {
				if (pan._fy_UIPanel.buttonsDown & (1 << button) != 0) {
					focusElement.Element.FyEMouseEvent(MouseEvent{
						Vec2i{0, 0},
						MOUSEEVENT_UP,
						button,
					})
				}
			}
			focusElement.Element.FyENormalEvent(FocusEvent{false})
		}
		pan._fy_UIPanel.focus = -1
	}
	pan._fy_UIPanel.content = content
}

func (pan *UIPanel) FyENormalEvent(ev NormalEvent) {
	if (pan._fy_UIPanel.focus != -1) {
		pan._fy_UIPanel.content[pan._fy_UIPanel.focus].Element.FyENormalEvent(ev)
	}
}

func (pan *UIPanel) _fy_Panel_ForwardMouseEvent(target PanelFixedElement, ev MouseEvent) {
	ev = ev.Offset(target.Pos.Negate())
	target.Element.FyEMouseEvent(ev)
}
func (pan *UIPanel) FyEMouseEvent(ev MouseEvent) {
	// Useful for debugging if any of the warnings come up
	//fmt.Printf("ui_core.go/Panel/FyEMouseEvent %v %v (%v, %v)\n", ev.Id, ev.Button, ev.Pos.X, ev.Pos.Y)
	invalid := false
	hittest := -1
	// Hit-test goes in reverse so that the element drawn last wins.
	for keyRev := range pan._fy_UIPanel.content {
		key := len(pan._fy_UIPanel.content) - (keyRev + 1)
		val := pan._fy_UIPanel.content[key]
		if (!val.Visible) {
			continue;
		}
		if (Area2iFromVecs(val.Pos, val.Element.FyESize()).Contains(ev.Pos)) {
			//fmt.Printf(" Hit index %v\n", key)
			hittest = key
			break
		}
	}
	switch (ev.Id) {
		case MOUSEEVENT_MOVE:
			// Mouse-move events go everywhere.
			for _, val := range pan._fy_UIPanel.content {
				pan._fy_Panel_ForwardMouseEvent(val, ev)
			}
			invalid = true
		case MOUSEEVENT_UP:
			if (pan._fy_UIPanel.buttonsDown & (1 << ev.Button) == 0) {
				fmt.Println("ui_core.go/Panel/FyEMouseEvent warning: Button removal on non-existent button")
				return
			} else {
				pan._fy_UIPanel.buttonsDown &= 0xFFFF ^ (1 << ev.Button)
			}
		case MOUSEEVENT_DOWN:
			if pan._fy_UIPanel.buttonsDown == 0 {
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
				if pan._fy_UIPanel.focus != hittest {
					// Note that this only happens when all other buttons have been released.
					// This prevents having to create fake release events.
					// The details of the order here are to do with issues when elements start modifying things in reaction to events.
					// Hence, the element that is being focused gets to run first so it will always receive an unfocus event after it has been focused.
					// While the element being unfocused is unlikely to get refocused under sane circumstances.
					// If worst comes to worst, make this stop sending focus events so nobody has to worry about focus state atomicity.
					oldFocus := pan._fy_UIPanel.focus
					pan._fy_UIPanel.focus = hittest
					newFocusFixed := PanelFixedElement{}
					if pan._fy_UIPanel.focus != -1 {
						newFocusFixed = pan._fy_UIPanel.content[pan._fy_UIPanel.focus]
					}
					// Since a mouse event came in in the first place, we know the panel's focused.
					// Focus the newly focused element.
					if newFocusFixed.Element != nil {
						newFocusFixed.Element.FyENormalEvent(FocusEvent{true})
					}
					// Unfocus the existing focused element, if any.
					if oldFocus != -1 {
						pan._fy_UIPanel.content[oldFocus].Element.FyENormalEvent(FocusEvent{false})
					}
				}
			}
			if (pan._fy_UIPanel.buttonsDown & (1 << ev.Button) != 0) {
				fmt.Println("ui_core.go/Panel/FyEMouseEvent warning: Button added when it was already added")
				return
			} else {
				pan._fy_UIPanel.buttonsDown |= (1 << ev.Button)
			}
	}
	// Yes, focus gets to receive mouse-move events out of bounds even if there are no buttons.
	// All the state is updated, forward the event
	if (!invalid && pan._fy_UIPanel.focus != -1) {
		pan._fy_Panel_ForwardMouseEvent(pan._fy_UIPanel.content[pan._fy_UIPanel.focus], ev)
	}
}
func (pan *UIPanel) FyEDraw(target Renderer, under bool) {
	if (pan.PanelClipping) {
		// Clipping: everything is inside panel bounds
		if (under) {
			return
		}
		oldClip := target.Clip()
		newClip := oldClip.Intersect(Area2iOfSize(pan.FyESize()))
		if (newClip.Empty()) {
			return
		}
		target.SetClip(newClip)
		defer target.SetClip(oldClip)
		for pass := 0; pass < 2; pass++ {
			for _, val := range pan._fy_UIPanel.content {
				if (!val.Visible) {
					continue;
				}
				target.Translate(val.Pos)
				val.Element.FyEDraw(target, pass == 0)
				target.Translate(val.Pos.Negate())
			}
		}
	} else {
		// Not clipping; this simply arranges a bunch of elements
		for _, val := range pan._fy_UIPanel.content {
			if (!val.Visible) {
				continue;
			}
			target.Translate(val.Pos)
			val.Element.FyEDraw(target, under)
			target.Translate(val.Pos.Negate())
		}
	}
}
func (pan *UIPanel) FyETick(f float64) {
	for _, val := range pan._fy_UIPanel.content {
		if (!val.Visible) {
			continue;
		}
		val.Element.FyETick(f)
	}
}

/*
 * "Proxy" element.
 */
type UIProxy struct {
	ProxyTarget UIElement
}
func (px *UIProxy) FyENormalEvent(ev NormalEvent) {
	px.ProxyTarget.FyENormalEvent(ev)
}
func (px *UIProxy) FyEMouseEvent(ev MouseEvent) {
	px.ProxyTarget.FyEMouseEvent(ev)
}
func (px *UIProxy) FyEDraw(target Renderer, under bool) {
	px.ProxyTarget.FyEDraw(target, under)
}
func (px *UIProxy) FyETick(f float64) {
	px.ProxyTarget.FyETick(f)
}
func (px *UIProxy) FyEResize(v Vec2i) {
	px.ProxyTarget.FyEResize(v)
}
func (px *UIProxy) FyESize() Vec2i {
	return px.ProxyTarget.FyESize()
}
