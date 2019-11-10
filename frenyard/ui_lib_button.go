package frenyard

// UIButtonThemedContent is a function to provide the visual components for a given UIButton state.
type UIButtonThemedContent func(hover bool, down bool) (NinePatchPackage, UILayoutElement)

// UIButton is a themable button.
type UIButton struct {
	UILayoutProxy
	_elem         UILayoutElement
	_innerOverlay *UIOverlayContainer
	_hover        bool
	_down         bool
	_theme        UIButtonThemedContent
	OnClick       func()
}

// NewUIButtonPtr creates a new UIButton.
func NewUIButtonPtr(theme UIButtonThemedContent, click func()) *UIButton {
	container := NewUIOverlayContainerPtr(NinePatchPackage{}, []UILayoutElement{})
	button := &UIButton{
		_innerOverlay: container,
		_theme:        theme,
		OnClick:       click,
	}
	InitUILayoutProxy(button, container)
	button._fyUIButtonUpdateState()
	return button
}

func (btn *UIButton) _fyUIButtonUpdateState() {
	npp, elem := btn._theme(btn._hover, btn._down)
	if elem != nil {
		btn._innerOverlay.SetContent(npp, []UILayoutElement{elem})
	} else {
		btn._innerOverlay.SetContent(npp, []UILayoutElement{})
	}
	btn._elem = elem
}

// FyEMouseEvent overrides UILayoutProxy.FyEMouseEvent
func (btn *UIButton) FyEMouseEvent(me MouseEvent) {
	lastHover := btn._hover
	lastDown := btn._down
	btn._hover = Area2iOfSize(btn.FyESize()).Contains(me.Pos)
	if me.Button == MouseButtonLeft {
		if me.ID == MouseEventUp {
			btn._down = false
			if btn._hover {
				btn.OnClick()
			}
		} else if me.ID == MouseEventDown {
			btn._down = true
		}
	}
	if lastHover != btn._hover || lastDown != btn._down {
		btn._fyUIButtonUpdateState()
	}
}
