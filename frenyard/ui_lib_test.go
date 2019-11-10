package frenyard

import "fmt"

/*
 * UIEventDebugger is a simple test element to debug events.
 */
type UIEventDebugger struct {
	UIElementComponent
	_mousePos Vec2i
}

// NewUIEventDebuggerPtr creates a UIEventDebugger
func NewUIEventDebuggerPtr(size Vec2i) *UIEventDebugger {
	return &UIEventDebugger{NewUIElementComponent(size), Vec2i{-8, -8}}
}

// FyENormalEvent implements UIElement.FyENormalEvent
func (cr *UIEventDebugger) FyENormalEvent(ev NormalEvent) {
	fmt.Printf("%+v\n", ev)
}

// FyEMouseEvent implements UIElement.FyEMouseEvent
func (cr *UIEventDebugger) FyEMouseEvent(ev MouseEvent) {
	fmt.Printf("%+v\n", ev)
	cr._mousePos = ev.Pos
}

// FyETick implements UIElement.FyETick
func (cr *UIEventDebugger) FyETick(deltaTime float64) {
}

// FyEDraw implements UIElement.FyEDraw
func (cr *UIEventDebugger) FyEDraw(target Renderer, under bool) {
	if !under {
		target.FillRect(0xFF3F51B5, Area2iOfSize(cr.FyESize()))
	} else {
		target.FillRect(0x40000000, Area2iOfSize(cr.FyESize()).Expand(Area2iFromVecs(Vec2i{-4, -4}, Vec2i{4, 4})))
		target.FillRect(0x40000000, Area2iOfSize(cr.FyESize()).Expand(Area2iFromVecs(Vec2i{-8, -8}, Vec2i{8, 8})))
	}
	var cursorRect Area2i
	if !under {
		cursorRect = Area2iOfSize(Vec2i{8, 8}).Translate(Vec2i{-4, -4})
	} else {
		cursorRect = Area2iOfSize(Vec2i{16, 16}).Translate(Vec2i{-8, -8})
	}
	target.Translate(cr._mousePos)
	target.FillRect(0x80FF0000, cursorRect)
	target.Translate(cr._mousePos.Negate())
}
