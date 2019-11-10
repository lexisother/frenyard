package frenyard

import (
	"runtime"
	"time"
	"github.com/veandco/go-sdl2/sdl"
)

type fySDL2Window struct {
	sdl2Renderer
	cacheMouse Vec2i
	window *sdl.Window
	receiver WindowReceiver
}

func (w *fySDL2Window) Name() string {
	{ z := sdl2Os(); defer z.End() }
	return w.window.GetTitle()
}
func (w *fySDL2Window) SetName(n string) {
	{ z := sdl2Os(); defer z.End() }
	w.window.SetTitle(n)
}

func (w *fySDL2Window) Destroy() {
	{ z := sdl2Os(); defer z.End() }
	w.base.osDelete()
	w.window.Destroy()
}

type fySDL2Backend struct {
	windows map[uint32]*fySDL2Window
}
var fyGlobalBackend *fySDL2Backend = &fySDL2Backend{map[uint32]*fySDL2Window{}}

func (r *fySDL2Backend) CreateWindow(name string, size Vec2i, vsync bool, receiver WindowReceiver) (Window, error) {
	{ z := sdl2Os(); defer z.End() }
	// Don't enable hidpi in SDL2. Pixels stop having a coherent meaning.
	if (vsync) {
		sdl.SetHint("SDL_HINT_RENDER_VSYNC", "1")
	} else {
		sdl.SetHint("SDL_HINT_RENDER_VSYNC", "0")
	}
	window, renderer, err := sdl.CreateWindowAndRenderer(size.X, size.Y, sdl.WINDOW_RESIZABLE)
	window.SetTitle(name)
	if err != nil {
		return nil, err
	}
	id, err := window.GetID()
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		return nil, err
	}
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	sWindow := &fySDL2Window{
		newSDL2Renderer(&sdl2RendererCore{renderer, false}),
		Vec2i{-1, -1},
		window,
		receiver,
	}
	fyGlobalBackend.windows[id] = sWindow
	receiver.FyRStart(sWindow)
	return sWindow, nil
}

func (r *fySDL2Backend) CreateTexture(size Vec2i, pixels []uint32) Texture {
	{ z := sdl2Os(); defer z.End() }
	if (len(pixels) != int(size.X) * int(size.Y)) {
		panic("invalid input to CreateTexture (size != pixels slice size)")
	}
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, size.X, size.Y,
 32, uint32(sdl.PIXELFORMAT_ARGB32))
	pixelsArray := surface.Pixels()
	for k, v := range pixels {
		pixelsArray[(k * 4) + 0] = byte((v & 0xFF000000) >> 24)
		pixelsArray[(k * 4) + 1] = byte((v & 0xFF0000) >> 16)
		pixelsArray[(k * 4) + 2] = byte((v & 0xFF00) >> 8)
		pixelsArray[(k * 4) + 3] = byte((v & 0xFF) >> 0)
	}
	if (err != nil) {
		panic(err)
	}
	return osSdl2SurfaceToFyTexture(surface)
}

func fySDL2MouseButton(button uint8) MouseButton {
	switch (button) {
		case sdl.BUTTON_LEFT:
			return MouseButtonLeft
		case sdl.BUTTON_MIDDLE:
			return MouseButtonMiddle
		case sdl.BUTTON_RIGHT:
			return MouseButtonRight
		case sdl.BUTTON_X1:
			return MouseButtonX1
		case sdl.BUTTON_X2:
			return MouseButtonX2
	}
	return MouseButtonNone
}
func fySDL2MouseWheelAdjuster(application WindowReceiver, cacheMouse Vec2i, apply int32, minus MouseButton, plus MouseButton) {
	for (apply < 0) {
		application.FyRMouseEvent(MouseEvent{cacheMouse, MouseEventDown, minus})
		application.FyRMouseEvent(MouseEvent{cacheMouse, MouseEventUp, minus})
		apply++;
	}
	for (apply > 0) {
		application.FyRMouseEvent(MouseEvent{cacheMouse, MouseEventDown, plus})
		application.FyRMouseEvent(MouseEvent{cacheMouse, MouseEventUp, plus})
		apply--;
	}
}

func init() {
	{ z := sdl2Os(); defer z.End() }
	sdl.Init(sdl.INIT_EVERYTHING)
	GlobalBackend = fyGlobalBackend
}

func (*fySDL2Backend) Run(ticker func(frameTime float64)) error {
	frameStart := time.Now()
	for !ExitFlag {
		for {
			timeLeft := TargetFrameTime - time.Since(frameStart).Seconds()
			if timeLeft <= 0 {
				break
			}
			time.Sleep(time.Duration(float64(time.Second) * timeLeft))
		}
		// Advance the frame officially
		lastFrameDuration := time.Since(frameStart)
		frameTime := lastFrameDuration.Seconds()
		frameStart = frameStart.Add(lastFrameDuration)

		// Unlock during sleep to avoid starving the OS thread.
		// But we don't want to thrash the OS thread when Renderer calls are going on,
		//  so calls back into the application are still on the OS thread.
		// If the application tries to run on multiple goroutines,
		//  the per-function OS thread locks will keep things from going off the rails.
		runtime.LockOSThread()
		ticker(frameTime)
		for _, wnd := range fyGlobalBackend.windows {
			wnd.receiver.FyRTick(frameTime)
		}
		fySDL2CRTCRegistry.osFlush()
		for {
			event := sdl.PollEvent()
			if (event == nil) {
				break
			}
			switch ev := event.(type) {
				case *sdl.MouseMotionEvent:
					window := fyGlobalBackend.windows[ev.WindowID]
					if (window != nil) {
						window.cacheMouse = Vec2i{ev.X, ev.Y}
						window.receiver.FyRMouseEvent(MouseEvent{window.cacheMouse, MouseEventMove, MouseButtonNone})
					}
				case *sdl.MouseButtonEvent:
					window := fyGlobalBackend.windows[ev.WindowID]
					if (window != nil) {
						window.cacheMouse = Vec2i{ev.X, ev.Y}
						buttonS := MouseEventDown
						if ev.State == sdl.RELEASED {
							buttonS = MouseEventUp
						}
						btn := fySDL2MouseButton(ev.Button)
						if btn != MouseButtonNone {
							window.receiver.FyRMouseEvent(MouseEvent{window.cacheMouse, buttonS, btn})
						}
					}
				case *sdl.MouseWheelEvent:
					window := fyGlobalBackend.windows[ev.WindowID]
					if (window != nil) {
						fySDL2MouseWheelAdjuster(window.receiver, window.cacheMouse, ev.X, MouseButtonScrollLeft, MouseButtonScrollRight)
						fySDL2MouseWheelAdjuster(window.receiver, window.cacheMouse, ev.Y, MouseButtonScrollUp, MouseButtonScrollDown)
					}
				case *sdl.WindowEvent:
					if ev.Event == sdl.WINDOWEVENT_CLOSE {
						window := fyGlobalBackend.windows[ev.WindowID]
						if window != nil {
							window.receiver.FyRClose()
						}
					}
				case *sdl.QuitEvent:
					// No way to determine cause, so for now won't even bother...
					ExitFlag = true
			}
		}
		runtime.UnlockOSThread()
		// OS thread
	}
	return nil
}
