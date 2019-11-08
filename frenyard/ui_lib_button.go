package frenyard

type UIButtonThemedContent func (hover bool, down bool) (NinePatchPackage, UILayoutElement)

type UIButton struct {
	UILayoutProxy
	_fy_UIButton_elem UILayoutElement
	_fy_UIButton_innerOverlay *UIOverlayContainer
	_fy_UIButton_hover bool
	_fy_UIButton_down bool
	_fy_UIButton_theme UIButtonThemedContent
	OnClick func ()
}
func NewUIButtonPtr(theme UIButtonThemedContent, click func ()) *UIButton {
	container := NewUIOverlayContainerPtr(NinePatchPackage{}, []UILayoutElement{})
	button := &UIButton{
		_fy_UIButton_innerOverlay: container,
		_fy_UIButton_theme: theme,
		OnClick: click,
	}
	button.ProxyTarget = container
	button._fy_UIButton_UpdateState()
	return button
}
func (btn *UIButton) _fy_UIButton_UpdateState() {
	npp, elem := btn._fy_UIButton_theme(btn._fy_UIButton_hover, btn._fy_UIButton_down)
	if elem != nil {
		btn._fy_UIButton_innerOverlay.SetContent(npp, []UILayoutElement{elem})
	} else {
		btn._fy_UIButton_innerOverlay.SetContent(npp, []UILayoutElement{})
	}
	btn._fy_UIButton_elem = elem
}
func (btn *UIButton) FyEMouseEvent(me MouseEvent) {
	lastHover := btn._fy_UIButton_hover
	lastDown := btn._fy_UIButton_down
	btn._fy_UIButton_hover = Area2iOfSize(btn.FyESize()).Contains(me.Pos)
	if me.Button == MOUSEBUTTON_LEFT {
		if me.Id == MOUSEEVENT_UP {
			btn._fy_UIButton_down = false
			if btn._fy_UIButton_hover {
				btn.OnClick()
			}
		} else if me.Id == MOUSEEVENT_DOWN {
			btn._fy_UIButton_down = true
		}
	}
	if lastHover != btn._fy_UIButton_hover || lastDown != btn._fy_UIButton_down {
		btn._fy_UIButton_UpdateState()
	}
}
