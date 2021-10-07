package framework

import (
	"fmt"
	"github.com/yellowsink/frenyard"
)

// UIEventDebugger is a simple test element to debug events.
type UIEventDebugger struct {
	UIElementComponent
	_mousePos frenyard.Vec2i
}

// NewUIEventDebuggerPtr creates a UIEventDebugger
func NewUIEventDebuggerPtr(size frenyard.Vec2i) *UIEventDebugger {
	return &UIEventDebugger{NewUIElementComponent(size), frenyard.Vec2i{-8, -8}}
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (cr *UIEventDebugger) FyENormalEvent(ev frenyard.NormalEvent) {
	fmt.Printf("%+v\n", ev)
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (cr *UIEventDebugger) FyEMouseEvent(ev frenyard.MouseEvent) {
	fmt.Printf("%+v\n", ev)
	cr._mousePos = ev.Pos
}

// FyETick implements UIElement.FyETick
func (cr *UIEventDebugger) FyETick(deltaTime float64) {
}

// FyEDraw implements UIElement.FyEDraw
func (cr *UIEventDebugger) FyEDraw(target frenyard.Renderer, under bool) {
	if !under {
		target.DrawRect(frenyard.DrawRectCommand{
			Colour: 0xFF3F51B5,
			Target: frenyard.Area2iOfSize(cr.FyESize()),
		})
	} else {
		target.DrawRect(frenyard.DrawRectCommand{
			Colour: 0x40000000,
			Target: frenyard.Area2iOfSize(cr.FyESize()).Expand(frenyard.Area2iFromVecs(frenyard.Vec2i{-4, -4}, frenyard.Vec2i{4, 4})),
		})
		target.DrawRect(frenyard.DrawRectCommand{
			Colour: 0x40000000,
			Target: frenyard.Area2iOfSize(cr.FyESize()).Expand(frenyard.Area2iFromVecs(frenyard.Vec2i{-8, -8}, frenyard.Vec2i{8, 8})),
		})
	}
	var cursorRect frenyard.Area2i
	if !under {
		cursorRect = frenyard.Area2iOfSize(frenyard.Vec2i{8, 8}).Translate(frenyard.Vec2i{-4, -4})
	} else {
		cursorRect = frenyard.Area2iOfSize(frenyard.Vec2i{16, 16}).Translate(frenyard.Vec2i{-8, -8})
	}
	target.Translate(cr._mousePos)
	target.DrawRect(frenyard.DrawRectCommand{
		Colour: 0x80FF0000,
		Target: cursorRect,
	})
	target.Translate(cr._mousePos.Negate())
}

// NewDebugWrapPtr wraps in an interior layout.
func NewDebugWrapPtr(inner UILayoutElement) UILayoutElement {
	return NewUIOverlayContainerPtr(nil, []UILayoutElement{ConvertElementToLayout(NewColouredRectPtr(0xFF200000, frenyard.Vec2i{})), inner})
}
