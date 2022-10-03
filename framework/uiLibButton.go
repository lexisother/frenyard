package framework

import (
	"github.com/uwu/frenyard"
)

// ButtonBehavior represents a button's behavior.
type ButtonBehavior func()

// UIButton is a themable button. Essentially, it's like an HTML button after the defaults are stripped off; the visual content is similar to any other element, but the "button" itself is an interaction proxy.
type UIButton struct {
	UILayoutProxy
	// "State" values
	Focused      bool
	Hover        bool
	Down         bool
	LastMousePos frenyard.Vec2i
	_behavior    func()
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
func (btn *UIButton) FyENormalEvent(me frenyard.NormalEvent) {
	btn.UILayoutProxy.FyENormalEvent(me)
	switch val := me.(type) {
	case FocusEvent:
		btn.Focused = val.Focused
	}
}

// FyEMouseEvent overrides UILayoutProxy.FyEMouseEvent
func (btn *UIButton) FyEMouseEvent(me frenyard.MouseEvent) {
	btn.UILayoutProxy.FyEMouseEvent(me)
	btn.Hover = frenyard.Area2iOfSize(btn.FyESize()).Contains(me.Pos)
	btn.LastMousePos = me.Pos
	if me.Button == frenyard.MouseButtonLeft {
		if me.ID == frenyard.MouseEventUp {
			btn.Down = false
			if btn.Hover {
				btn._behavior()
			}
		} else if me.ID == frenyard.MouseEventDown {
			btn.Down = true
		}
	}
}
