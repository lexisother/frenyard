package frenyard

// ButtonBehavior represents a button's behavior.
type ButtonBehavior func()

// UIButton is a themable button. Essentially, it's like an HTML button after the defaults are stripped off; the visual content is similar to any other element, but the "button" itself is an interaction proxy.
type UIButton struct {
	UILayoutProxy
	// "State" values
	Focused       bool
	Hover         bool
	Down          bool
	_behavior     func()
}

// NewUIButtonPtr creates a new UIButton.
func NewUIButtonPtr(theme UILayoutElement, click ButtonBehavior) *UIButton {
	button := &UIButton{
		_behavior: click,
	}
	InitUILayoutProxy(button, theme)
	return button
}

// FyENormalEvent overrides UILayoutProxy.FyENormalEvent
func (btn *UIButton) FyENormalEvent(me NormalEvent) {
	btn.UILayoutProxy.FyENormalEvent(me)
	switch val := me.(type) {
		case FocusEvent:
			btn.Focused = val.Focused
	}
}

// FyEMouseEvent overrides UILayoutProxy.FyEMouseEvent
func (btn *UIButton) FyEMouseEvent(me MouseEvent) {
	btn.UILayoutProxy.FyEMouseEvent(me)
	btn.Hover = Area2iOfSize(btn.FyESize()).Contains(me.Pos)
	if me.Button == MouseButtonLeft {
		if me.ID == MouseEventUp {
			btn.Down = false
			if btn.Hover {
				btn._behavior()
			}
		} else if me.ID == MouseEventDown {
			btn.Down = true
		}
	}
}
