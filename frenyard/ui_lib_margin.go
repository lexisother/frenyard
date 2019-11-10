package frenyard

// NewUIMarginContainerPtr is a wrapper around NewUIOverlayContainerPtr which creates a margin without any other details around a given element.
func NewUIMarginContainerPtr(innards UILayoutElement, margin Area2i) UILayoutElement {
	return NewUIOverlayContainerPtr(NinePatchPackage{
		Padding: margin,
	}, []UILayoutElement{
		innards,
	})
}
