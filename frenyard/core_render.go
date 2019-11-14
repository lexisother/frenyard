/*
 * PLEASE KEEP IN MIND: NONE OF THIS IS 'PUBLIC' FOR DEPENDENCY PURPOSES (yet)
 */

package frenyard

// ExitFlag when set to true, exits the application.
var ExitFlag bool = false

// Backend is the set of "entrypoint" functions to the core API.
type Backend interface {
	// Begins the frame loop. Stops when ExitFlag is set to true.
	Run(ticker func(frameTime float64)) error
	CreateWindow(name string, size Vec2i, vsync bool, receiver WindowReceiver) (Window, error)
	CreateTexture(size Vec2i, pixels []uint32) Texture
}

// GlobalBackend is the global instance of Backend.
var GlobalBackend Backend

// TargetFrameTime controls the framerate of the application.
var TargetFrameTime float64 = 0.05 // 20FPS

// WindowReceiver receives window events.
type WindowReceiver interface {
	FyRStart(w Window)
	FyRTick(time float64)
	FyRNormalEvent(n NormalEvent)
	FyRMouseEvent(m MouseEvent)
	// Note: The window is destroyed after this completes.
	FyRClose()
}

// Window type. This type MAY be user-implemented with the understanding that no Core/CoreExt functions accept Window or Renderer (hence there are no potential issues), and that this is still intended as a Core/CoreExt type so doesn't get the name-prefixing.
type Window interface {
	Renderer
	Name() string
	SetName(name string)
	Present()
	Destroy()
	// Gets the DPI of the window. This can change. Oh well.
	GetLocalDPI() float64
}

// Renderer is an abstract rendering interface.
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

// Texture interface. This is automatically deleted on finalization.
type Texture interface {
	Size() Vec2i
}

/*
 * There are two kinds of event: Normal (focus-target) events and Mouse Events.
 * Mouse Events get lots of special treatment as they bypass focus targeting, need offset logic, etc.
 * Normal Events have the same handling no matter what, so they don't even need to be considered in core.
 * There's some UI special handling for Focus/Unfocus events, but that's defined out of core anyway
 */

// NormalEvent is the base of event types that target the focused element.
type NormalEvent interface {
}

// MouseEventID describes a type of MouseEvent.
type MouseEventID uint8

// MouseEventMove indicates that the event is because the mouse was moved.
const MouseEventMove MouseEventID = 0

// MouseEventDown indicates that the event is because a mouse button was pressed
const MouseEventDown MouseEventID = 1

// MouseEventUp indicates that the event is because a mouse button was released
const MouseEventUp MouseEventID = 2

// MouseButton describes a mouse button.
type MouseButton int8

// MouseButtonNone indicates that no button was involved. Only appears for MouseEventMove
const MouseButtonNone MouseButton = -1

// Numbers from here forward must be allocatable in _fy_Panel_ButtonsDown
// Including the scroll buttons!

// MouseButtonLeft is the left mouse button.
const MouseButtonLeft MouseButton = 0

// MouseButtonMiddle is the middle mouse button (Do be warned: Laptop users do not get this in any 'easy to understand' form.)
const MouseButtonMiddle MouseButton = 1

// MouseButtonRight is the right mouse button
const MouseButtonRight MouseButton = 2

// MouseButtonX1 is a fancy auxiliary mouse button that not all people have
const MouseButtonX1 MouseButton = 3

// MouseButtonX2 is a fancy auxiliary mouse button that not all people have
const MouseButtonX2 MouseButton = 4

// These are a form of button because it simplifies the implementation massively at no real cost.

// MouseButtonScrollUp is a virtual scroll
const MouseButtonScrollUp MouseButton = 5

// MouseButtonScrollDown is a virtual scroll
const MouseButtonScrollDown MouseButton = 6

// MouseButtonScrollLeft is a virtual scroll
const MouseButtonScrollLeft MouseButton = 7

// MouseButtonScrollRight is a virtual scroll
const MouseButtonScrollRight MouseButton = 8

// MouseButtonLength is not a real button. You may need to use (int8)(0) in for loops.
const MouseButtonLength MouseButton = 9

// MouseEvent is a mouse event.
type MouseEvent struct {
	// Where the mouse is *relative to the receiving element.*
	Pos Vec2i
	// Indicates the sub-type of the event.
	ID MouseEventID
	// Meaningless for MOUSEEVENT_MOVE. See MOUSEBUTTON_*
	Button MouseButton
}

// Offset offsets the mouse event by a given amount.
func (ev MouseEvent) Offset(offset Vec2i) MouseEvent {
	return MouseEvent{
		Pos:    ev.Pos.Add(offset),
		ID:     ev.ID,
		Button: ev.Button,
	}
}
