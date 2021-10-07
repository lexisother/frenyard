package framework

import "github.com/yellowsink/frenyard"

// SlideTransition represents a specific transition queued for a UISlideTransitionContainer.
type SlideTransition struct {
	Element  UILayoutElement
	Length   float64
	Reverse  bool
	Vertical bool
}

// UISlideTransitionContainer creates a 'slide transition' box.
type UISlideTransitionContainer struct {
	UIPanel
	UILayoutElementComponent
	// The current transition, with the current element.
	_transition SlideTransition
	// If != nil, then transition is active and this is the element moving out.
	_last UILayoutElement
	// The transition queue. When a transition ends (and only then), this queue is checked.
	_transitionQueue []SlideTransition
	// Timer for the transition (from 0 to _transitionLength)
	_transitionTime float64
}

// NewUISlideTransitionContainerPtr creates a 'slide transition' box.
func NewUISlideTransitionContainerPtr(initContent UILayoutElement) *UISlideTransitionContainer {
	container := &UISlideTransitionContainer{
		UIPanel: NewPanel(frenyard.Vec2i{}),
		_transition: SlideTransition{
			Element: initContent,
			Length:  1.0,
		},
		_transitionQueue: []SlideTransition{},
		_transitionTime:  1.0,
	}
	InitUILayoutElementComponent(container)
	if initContent != nil {
		container.ThisUILayoutElementComponentDetails.Attach(initContent)
		container.FyEResize(initContent.FyESize())
		container._fyUpdatePositions()
	}
	return container
}

// FyLSubelementChanged implements UILayoutElement.FyLSubelementChanged
func (ufc *UISlideTransitionContainer) FyLSubelementChanged() {
	ufc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (ufc *UISlideTransitionContainer) FyLSizeForLimits(limits frenyard.Vec2i) frenyard.Vec2i {
	if ufc._transition.Element != nil {
		return ufc._transition.Element.FyLSizeForLimits(limits)
	}
	return frenyard.Vec2i{}
}

// This is the core 'start a transition' function. It overwrites the previous transition.
func (ufc *UISlideTransitionContainer) setTransition(next SlideTransition) {
	if ufc._last != nil {
		ufc.ThisUILayoutElementComponentDetails.Detach(ufc._last)
	}
	ufc._last = ufc._transition.Element
	if next.Length > 0.0 {
		// Completion takes non-zero time.
		ufc._transitionTime = 0
	} else {
		// Cause instant completion.
		ufc._transitionTime = 1
		next.Length = 1
	}
	ufc._transition = next
	ufc.ThisUILayoutElementComponentDetails.Attach(ufc._transition.Element)
	ufc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// TransitionTo transitions to a new slide.
func (ufc *UISlideTransitionContainer) TransitionTo(next SlideTransition) {
	if ufc._last != nil {
		ufc._transitionQueue = append(ufc._transitionQueue, next)
	} else {
		ufc.setTransition(next)
	}
}

// FyEResize overrides UIPanel.FyEResize
func (ufc *UISlideTransitionContainer) FyEResize(size frenyard.Vec2i) {
	ufc.UIPanel.FyEResize(size)
	if ufc._transition.Element != nil {
		ufc._transition.Element.FyEResize(size)
	}
	if ufc._last != nil {
		ufc._last.FyEResize(size)
	}
	ufc._fyUpdatePositions()
}
func (ufc *UISlideTransitionContainer) _fyUpdatePositions() {
	areaSize := ufc.FyESize().ConditionalTranspose(ufc._transition.Vertical).X
	transitionIP := ufc._transitionTime / ufc._transition.Length
	transitionIP = frenyard.EasingInOut(frenyard.EasingQuadraticIn)(transitionIP)
	point := int32(float64(areaSize) * transitionIP)
	mul := int32(1)
	if ufc._transition.Reverse {
		mul *= -1
	}
	feA := PanelFixedElement{
		Pos:     frenyard.Vec2i{(-point) * mul, 0}.ConditionalTranspose(ufc._transition.Vertical),
		Element: ufc._last,
		Locked:  true,
		Visible: true,
	}
	feB := PanelFixedElement{
		Pos:     frenyard.Vec2i{(areaSize - point) * mul, 0}.ConditionalTranspose(ufc._transition.Vertical),
		Element: ufc._transition.Element,
		Visible: true,
		Locked:  ufc._last != nil,
	}
	var slice []PanelFixedElement
	if ufc._transition.Element != nil {
		if ufc._last != nil {
			slice = []PanelFixedElement{feA, feB}
		} else {
			feB.Pos = frenyard.Vec2i{}
			slice = []PanelFixedElement{feB}
		}
	} else {
		if ufc._last != nil {
			slice = []PanelFixedElement{feA}
		} else {
			slice = []PanelFixedElement{}
		}
	}
	ufc.ThisUIPanelDetails.SetContent(slice)
}

// FyETick overrides UIPanel.FyETick
func (ufc *UISlideTransitionContainer) FyETick(time float64) {
	ufc.UIPanel.FyETick(time)
	if ufc._last != nil {
		ufc._transitionTime += time
		if ufc._transitionTime >= ufc._transition.Length {
			ufc._transitionTime = ufc._transition.Length
			ufc.ThisUILayoutElementComponentDetails.Detach(ufc._last)
			ufc._last = nil
			// Ok, completed a transition. Check if there's another one incoming
			if len(ufc._transitionQueue) > 0 {
				// Yup. Begin & remove the first transition.
				ufc.setTransition(ufc._transitionQueue[0])
				ufc._transitionQueue = ufc._transitionQueue[1:]
			}
		}
		ufc._fyUpdatePositions()
	}
}
