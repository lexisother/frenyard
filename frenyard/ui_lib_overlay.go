package frenyard

// 'Overlay': Overlays elements over each other.
// Can also install a nine-patch-alike.
// Useful for backgrounds.

type NinePatchPackage struct {
	Over NinePatch
	Under NinePatch
	Padding Area2i
	Clipping bool
}

type UIOverlayContainer struct {
	UIPanel
	UILayoutElementComponent
	_fy_UIOverlayContainer_NinePatch NinePatchPackage
	_fy_UIOverlayContainer_State []UILayoutElement
	_fy_UIOverlayContainer_PreferredSize Vec2i
}

func NewUIOverlayContainerPtr(npp NinePatchPackage, setup []UILayoutElement) *UIOverlayContainer {
	container := &UIOverlayContainer{
		UIPanel: NewPanel(Vec2i{}),
	}
	InitUILayoutElement(container)
	container.SetContent(npp, setup)
	container.FyEResize(container._fy_UIOverlayContainer_PreferredSize)
	return container
}

func (ufc *UIOverlayContainer) FyLSubelementChanged() {
	size := Vec2i{}
	for _, v := range ufc._fy_UIOverlayContainer_State {
		size = size.Max(v.FyLSizeForLimits(Vec2iUnlimited()))
	}
	ufc._fy_UIOverlayContainer_PreferredSize = size.Add(ufc._fy_UIOverlayContainer_NinePatch.Padding.Size())
	ufc.UIThis.ContentChanged()
}

func (ufc *UIOverlayContainer) FyLSizeForLimits(limits Vec2i) Vec2i {
	if limits.Ge(ufc._fy_UIOverlayContainer_PreferredSize) {
		return ufc._fy_UIOverlayContainer_PreferredSize
	}
	max := Vec2i{}
	for _, v := range ufc._fy_UIOverlayContainer_State {
		max = max.Max(v.FyLSizeForLimits(limits.Add(ufc._fy_UIOverlayContainer_NinePatch.Padding.Size().Negate())))
	}
	return max.Add(ufc._fy_UIOverlayContainer_NinePatch.Padding.Size())
}

func (ufc *UIOverlayContainer) SetContent(npp NinePatchPackage, setup []UILayoutElement) {
	if ufc._fy_UIOverlayContainer_State != nil {
		for _, v := range ufc._fy_UIOverlayContainer_State {
			ufc.UIThis.Detach(v)
		}
	}
	ufc._fy_UIOverlayContainer_State = setup
	ufc._fy_UIOverlayContainer_NinePatch = npp
	ufc.PanelClipping = npp.Clipping
	for _, v := range setup {
		ufc.UIThis.Attach(v)
	}
	ufc.FyLSubelementChanged()
}

func (ufc *UIOverlayContainer) FyEResize(size Vec2i) {
	ufc.UIElementComponent.FyEResize(size)
	area := Area2iOfSize(size).Contract(ufc._fy_UIOverlayContainer_NinePatch.Padding)
	fixes := make([]PanelFixedElement, len(ufc._fy_UIOverlayContainer_State))
	for idx, slot := range ufc._fy_UIOverlayContainer_State {
		fixes[idx] = PanelFixedElement{
			Pos: area.Pos(),
			Visible: true,
			Element: slot,
		}
		slot.FyEResize(area.Size())
	}
	ufc.UIPanel.UIPanelSetContent(fixes)
}

func (ufc *UIOverlayContainer) FyEDraw(r Renderer, under bool) {
	if under {
		ufc._fy_UIOverlayContainer_NinePatch.Under.Draw(r, Area2iOfSize(ufc.FyESize()))
	} else {
		ufc._fy_UIOverlayContainer_NinePatch.Over.Draw(r, Area2iOfSize(ufc.FyESize()))
	}
	ufc.UIPanel.FyEDraw(r, under)
}
