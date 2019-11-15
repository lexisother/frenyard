package frenyard

// 'Overlay': Overlays elements over each other.
// Can also install a nine-patch-alike.
// Useful for backgrounds.

// NinePatchPackage packages several NinePatches into a set that can be drawn in a UI context.
type NinePatchPackage struct {
	Over     NinePatch
	Under    NinePatch
	Padding  Area2i
	// Scales everything (including padding!)
	Scale    float64
	Clipping bool
}

// GetEffectivePadding scales Padding by the scale, which provides the padding as it is used in practice.
func (npp NinePatchPackage) GetEffectivePadding() Area2i {
	return ScaleMargin2(npp.Scale, npp.Padding, ScaleRMNinePatch)
}

// UIOverlayContainer overlays elements on top of each other, and this itself on top of a potentially padded NinePatchPackage.
type UIOverlayContainer struct {
	UIPanel
	UILayoutElementComponent
	_ninePatch     NinePatchPackage
	_state         []UILayoutElement
	_preferredSize Vec2i
}

// NewUIOverlayContainerPtr creates a UIOverlayContainer
func NewUIOverlayContainerPtr(npp NinePatchPackage, setup []UILayoutElement) *UIOverlayContainer {
	container := &UIOverlayContainer{
		UIPanel: NewPanel(Vec2i{}),
	}
	InitUILayoutElementComponent(container)
	container.SetContent(npp, setup)
	container.FyEResize(container._preferredSize)
	return container
}

// FyLSubelementChanged implements UILayoutElement.FyLSubelementChanged
func (ufc *UIOverlayContainer) FyLSubelementChanged() {
	size := Vec2i{}
	for _, v := range ufc._state {
		size = size.Max(v.FyLSizeForLimits(Vec2iUnlimited()))
	}
	ufc._preferredSize = size.Add(ufc._ninePatch.GetEffectivePadding().Size())
	ufc.ThisUILayoutElementComponentDetails.ContentChanged()
}

// FyLSizeForLimits implements UILayoutElement.FyLSizeForLimits
func (ufc *UIOverlayContainer) FyLSizeForLimits(limits Vec2i) Vec2i {
	if limits.Ge(ufc._preferredSize) {
		return ufc._preferredSize
	}
	max := Vec2i{}
	paddingSize := ufc._ninePatch.GetEffectivePadding().Size()
	for _, v := range ufc._state {
		max = max.Max(v.FyLSizeForLimits(limits.Add(paddingSize.Negate())))
	}
	return max.Add(paddingSize)
}

// SetContent changes the content of the UIOverlayContainer.
func (ufc *UIOverlayContainer) SetContent(npp NinePatchPackage, setup []UILayoutElement) {
	if ufc._state != nil {
		for _, v := range ufc._state {
			ufc.ThisUILayoutElementComponentDetails.Detach(v)
		}
	}
	ufc._state = setup
	ufc._ninePatch = npp
	ufc.ThisUIPanelDetails.Clipping = npp.Clipping
	for _, v := range setup {
		ufc.ThisUILayoutElementComponentDetails.Attach(v)
	}
	ufc.FyLSubelementChanged()
}

// FyEResize overrides UIPanel.FyEResize
func (ufc *UIOverlayContainer) FyEResize(size Vec2i) {
	ufc.UIPanel.FyEResize(size)
	area := Area2iOfSize(size).Contract(ufc._ninePatch.GetEffectivePadding())
	fixes := make([]PanelFixedElement, len(ufc._state))
	for idx, slot := range ufc._state {
		fixes[idx] = PanelFixedElement{
			Pos:     area.Pos(),
			Visible: true,
			Element: slot,
		}
		slot.FyEResize(area.Size())
	}
	ufc.ThisUIPanelDetails.SetContent(fixes)
}

// FyEDraw overrides UIPanel.FyEDraw
func (ufc *UIOverlayContainer) FyEDraw(r Renderer, under bool) {
	if under {
		ufc._ninePatch.Under.Draw(r, Area2iOfSize(ufc.FyESize()), ufc._ninePatch.Scale)
	} else {
		ufc._ninePatch.Over.Draw(r, Area2iOfSize(ufc.FyESize()), ufc._ninePatch.Scale)
	}
	ufc.UIPanel.FyEDraw(r, under)
}
