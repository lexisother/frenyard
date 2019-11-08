package frenyard

import "fmt"

/*
 * UIEventDebugger is a simple test element to debug events.
 */
type UIEventDebugger struct {
	UIElementComponent
	_fy_EventDebugger_MousePos Vec2i
}
func NewUIEventDebugger(size Vec2i) UIEventDebugger {
	return UIEventDebugger{NewUIElementComponent(size), Vec2i{-8, -8}}
}
func (cr *UIEventDebugger) FyENormalEvent(ev NormalEvent) {
	fmt.Printf("%+v\n", ev)
}
func (cr *UIEventDebugger) FyEMouseEvent(ev MouseEvent) {
	fmt.Printf("%+v\n", ev)
	cr._fy_EventDebugger_MousePos = ev.Pos
}
func (cr *UIEventDebugger) FyETick(deltaTime float64) {
}
func (cr *UIEventDebugger) FyEDraw(target Renderer, under bool) {
	if (!under) {
		target.FillRect(0xFF3F51B5, Area2iOfSize(cr.FyESize()))
	} else {
		target.FillRect(0x40000000, Area2iOfSize(cr.FyESize()).Expand(Area2iFromVecs(Vec2i{-4, -4}, Vec2i{4, 4})))
		target.FillRect(0x40000000, Area2iOfSize(cr.FyESize()).Expand(Area2iFromVecs(Vec2i{-8, -8}, Vec2i{8, 8})))
	}
	var cursorRect Area2i
	if (!under) {
		cursorRect = Area2iOfSize(Vec2i{8, 8}).Translate(Vec2i{-4, -4})
	} else {
		cursorRect = Area2iOfSize(Vec2i{16, 16}).Translate(Vec2i{-8, -8})
	}
	target.Translate(cr._fy_EventDebugger_MousePos)
	target.FillRect(0x80FF0000, cursorRect)
	target.Translate(cr._fy_EventDebugger_MousePos.Negate())
}
