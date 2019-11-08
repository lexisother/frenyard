package frenyard

func NewUIMarginContainerPtr(innards UILayoutElement, margin Area2i) UILayoutElement {
	return NewUIOverlayContainerPtr(NinePatchPackage{
		Padding: margin,
	}, []UILayoutElement{
		innards,
	})
}
