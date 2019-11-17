package frenyard

// UISlideTransitionContainer creates a 'slide transition' box.
type UISlideTransitionContainer struct {
	UIPanel
	UILayoutElementComponent
	// The main (or moving-in) element. 
	_main UILayoutElement
	// If != nil, then transition is active and this is the element moving out.
	_last UILayoutElement
	// Timer for the transition (from 0 to _transitionLength)
	_transitionTime float64
	// The amount of time the transition takes.
	_transitionLength float64
	// If the transition goes in reverse.
	_transitionReverse bool
	// If the transition is vertical.
	_transitionVertical bool
}

// NewUISlideTransitionContainerPtr creates a 'slide transition' box.
func NewUISlideTransitionContainerPtr(initContent UILayoutElement) *UISlideTransitionContainer {
	container := &UISlideTransitionContainer{
		UIPanel: NewPanel(Vec2i{}),
		_main: initContent,
		_transitionTime: 1.0,
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
func (ufc *UISlideTransitionContainer) FyLSizeForLimits(limits Vec2i) Vec2i {
	if ufc._main != nil {
		return ufc._main.FyLSizeForLimits(limits)
	}
	return Vec2i{}
}

// TransitionTo transitions to a new slide.
func (ufc *UISlideTransitionContainer) TransitionTo(next UILayoutElement, time float64, reverse bool, vertical bool) {
	if ufc._last != nil {
		ufc.ThisUILayoutElementComponentDetails.Detach(ufc._last)
	}
	ufc._last = ufc._main
	ufc._main = next
	ufc.ThisUILayoutElementComponentDetails.Attach(ufc._main)
	if time > 0.0 {
		// Completion takes non-zero time.
		ufc._transitionTime = 0
		ufc._transitionLength = time
	} else {
		// Cause instant completion.
		ufc._transitionTime = 1
		ufc._transitionLength = 1
	}
	ufc._transitionReverse = reverse
	ufc._transitionVertical = vertical
	ufc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// FyEResize overrides UIPanel.FyEResize
func (ufc *UISlideTransitionContainer) FyEResize(size Vec2i) {
	ufc.UIPanel.FyEResize(size)
	if ufc._main != nil {
		ufc._main.FyEResize(size)
	}
	if ufc._last != nil {
		ufc._last.FyEResize(size)
	}
	ufc._fyUpdatePositions()
}
func (ufc *UISlideTransitionContainer) _fyUpdatePositions() {
	areaSize := ufc.FyESize().ConditionalTranspose(ufc._transitionVertical).X
	transitionIP := ufc._transitionTime / ufc._transitionLength
	transitionIP = EasingInOut(EasingQuadraticIn)(transitionIP)
	point := int32(float64(areaSize) * transitionIP)
	mul := int32(1)
	if ufc._transitionReverse {
		mul *= -1
	}
	feA := PanelFixedElement{
		Pos: Vec2i{(-point) * mul, 0}.ConditionalTranspose(ufc._transitionVertical),
		Element: ufc._last,
		Locked: true,
		Visible: true,
	}
	feB := PanelFixedElement{
		Pos: Vec2i{(areaSize - point) * mul, 0}.ConditionalTranspose(ufc._transitionVertical),
		Element: ufc._main,
		Visible: true,
	}
	var slice []PanelFixedElement
	if ufc._main != nil {
		if ufc._last != nil {
			slice = []PanelFixedElement{feA, feB}
		} else {
			feB.Pos = Vec2i{}
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
		if ufc._transitionTime >= ufc._transitionLength {
			ufc._transitionTime = ufc._transitionLength
			ufc.ThisUILayoutElementComponentDetails.Detach(ufc._last)
			ufc._last = nil
		}
		ufc._fyUpdatePositions()
	}
}
