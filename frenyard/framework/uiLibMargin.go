package framework

import (
	"github.com/20kdc/CCUpdaterUI/frenyard"
)

// NewUIMarginContainerPtr is a wrapper around NewUIOverlayContainerPtr which creates a margin without any other details around a given element.
func NewUIMarginContainerPtr(innards UILayoutElement, margin frenyard.Area2i) UILayoutElement {
	return NewUIOverlayContainerPtr(&NinePatchFrame{
		Padding: margin,
	}, []UILayoutElement{
		innards,
	})
}
