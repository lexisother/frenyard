/*
 * PLEASE KEEP IN MIND: NONE OF THIS IS 'PUBLIC' FOR DEPENDENCY PURPOSES (yet)
 */

package frenyard

var ExitFlag bool = false

/*
 * Access via frenyard.GlobalBackend defined below
 */
type Backend interface {
	// Begins the frame loop. Stops when ExitFlag is set to true.
	Run(ticker func(frameTime float64)) error
	CreateWindow(name string, size Vec2i, vsync bool, receiver WindowReceiver) (Window, error)
	CreateTexture(size Vec2i, pixels []uint32) Texture
}
var GlobalBackend Backend
var TargetFrameTime float64 = 0.05 // 20FPS

type WindowReceiver interface {
	FyRStart(w Window)
	FyRTick(time float64)
	FyRNormalEvent(n NormalEvent)
	FyRMouseEvent(m MouseEvent)
	// Note: "Closed" does not mean "Destroyed".
	FyRClose()
}

/*
 * This type MAY be user-implemented with the understanding that:
 * 1. No Core/CoreExt functions accept Window or Renderer.
 * As such there is no such thing as a "fake" Window or Renderer.
 * 2. The conflict-prevention name prefixes don't apply.
 */
type Window interface {
	Renderer
	Name() string
	SetName(name string)
	Present()
	Destroy()
}

/*
 * Abstract rendering interface.
 * Provided by an 'application start' function such as RunSDL2.
 * Some things will get refactored when this has to inevitably go multi-window.
 * Also note that renderers may be translated and scissored.
 */
type Renderer interface {	
	/* Draw */

	// Fills a rectangle.
	FillRect(colour uint32, target Area2i)
	// Stretches a texture.
	TexRect(sheet Texture, colour uint32, sprite Area2i, target Area2i)
	
	/* Control */
	
	// Translates the renderer's target. Use with defer to undo later.
	Translate(vec Vec2i)
	// Sets the clip area, relative to the current translation.
	SetClip(clip Area2i)
	// Gets the clip area, relative to the current translation.
	Clip() Area2i
	
	/* Outer Control */
	
	// Gets the size of the drawing area.
	Size() Vec2i
	// Clears & resets clip/translate/etc. You should do this at the start of  frame.
	Reset(colour uint32)
}

/*
 * Texture interface. This is automatically deleted on finalization.
 */
type Texture interface {
	Size() Vec2i
}

/*
 * There are two kinds of event: Normal (focus-target) events and Mouse Events.
 * Mouse Events get lots of special treatment as they bypass focus targetting, need offset logic, etc.
 * Normal Events have the same handling no matter what, so they don't even need to be considered in core.
 * There's some UI special handling for Focus/Unfocus events, but that's defined out of core anyway
 */
type NormalEvent interface {
}

const MOUSEEVENT_MOVE uint8 = 0
const MOUSEEVENT_DOWN uint8 = 1
const MOUSEEVENT_UP uint8 = 2

const MOUSEBUTTON_NONE int8 = -1
// Numbers from here forward must be allocatable in _fy_Panel_ButtonsDown
// Including the scroll buttons!
const MOUSEBUTTON_LEFT int8 = 0
const MOUSEBUTTON_MIDDLE int8 = 1
const MOUSEBUTTON_RIGHT int8 = 2
const MOUSEBUTTON_X1 int8 = 3
const MOUSEBUTTON_X2 int8 = 4
// These are a form of button because it simplifies the implementation massively at no real cost.
const MOUSEBUTTON_SCROLL_UP int8 = 5
const MOUSEBUTTON_SCROLL_DOWN int8 = 6 
const MOUSEBUTTON_SCROLL_LEFT int8 = 7
const MOUSEBUTTON_SCROLL_RIGHT int8 = 8
// Not a real button. You may need to use (int8)(0) in for loops
const MOUSEBUTTON_LENGTH int8 = 9

/*
 * Mouse event.
 */
type MouseEvent struct {
	// Where the mouse is *relative to the receiving element.*
	Pos Vec2i
	// Indicates the sub-type of the event.
	Id uint8
	// Meaningless for MOUSEEVENT_MOVE. See MOUSEBUTTON_*
	Button int8
}

func (ev MouseEvent) Offset(offset Vec2i) MouseEvent {
	return MouseEvent{
		Pos: ev.Pos.Add(offset),
		Id: ev.Id,
		Button: ev.Button,
	}
}
