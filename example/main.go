package main

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/example/screens"
	"github.com/lexisother/frenyard/example/src"
	"github.com/lexisother/frenyard/framework"
)

func main() {
	// This controls the framerate
	frenyard.TargetFrameTime = 0.016

	// Initialize our main container capable of slide animations
	slideContainer := framework.NewUISlideTransitionContainerPtr(nil)
	slideContainer.FyEResize(design.SizeWindowInit)

	// Initialize our window and bind it to our container element
	wnd, err := framework.CreateBoundWindow("Frenyard UI Playground", true, design.ThemeBackground, slideContainer)
	if err != nil {
		panic(err)
	}

	// Setup sets the sizes, fonts and borders according to a reasonable inferred
	// scale.
	design.Setup(frenyard.InferScale(wnd))
	wnd.SetSize(design.SizeWindow)

	// Initialize our app struct, referencable from anywhere.
	app := &src.UpApplication{
		MainContainer:    slideContainer,
		Window:           wnd,
		UpQueued:         make(chan func(), 16),
		TeleportSettings: framework.SlideTransition{},
	}

	screens.SetupButtons()
	screens.SetupTwo()

	app.ShowPrimaryView()

	// Start the backend. Stops when ExitFlag is set to true.
	frenyard.GlobalBackend.Run(func(frameTime float64) {
		select {
		case fn := <-app.UpQueued:
			fn()
		default:
		}
	})
}
