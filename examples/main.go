package main

import (
	"github.com/uwu/frenyard"
	"github.com/uwu/frenyard/design"
	"github.com/uwu/frenyard/framework"
)

func main() {
	// This controls the framerate
	frenyard.TargetFrameTime = 0.016

	// Initialize our main container capable of slide animations
	slideContainer := framework.NewUISlideTransitionContainerPtr(nil)
	slideContainer.FyEResize(design.SizeWindowInit)

	// Initialize our window and bind it to our container element
	wnd, err := framework.CreateBoundWindow("Frenyard Example", true, design.ThemeBackground, slideContainer)
	if err != nil {
		panic(err)
	}

	// Setup sets the sizes, fonts and borders according to a reasonable inferred
	// scale.
	design.Setup(frenyard.InferScale(wnd))
	wnd.SetSize(design.SizeWindow)

	// Initialize our app struct, referencable from anywhere.
	app := (&UpApplication{
		MainContainer:    slideContainer,
		Window:           wnd,
		UpQueued:         make(chan func(), 16),
		TeleportSettings: framework.SlideTransition{},
	})

	strptr := ""

	// Start an instant transition to our main screen.
	app.Teleport(
		// A 'document', with a title header and body.
		design.LayoutDocument(
			design.Header{
				Title: "Example App",
			},
			// Main flexbox container, contains all elements.
			framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
				DirVertical: true,
				Slots: []framework.FlexboxSlot{
					{
						Element: //framework.NewUILabelPtr(integration.NewTextTypeChunk("Hello World!", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
							design.NewUITextboxPtr("test", &strptr, "hi!"),
					},
					{
						Element: design.NewUITextboxPtr("", &strptr, "hi 2!"),
					},
				},
			}),
			// Sets the flexbox container to be scrollable.
			true,
		),
	)

	// Start the backend. Stops when ExitFlag is set to true.
	frenyard.GlobalBackend.Run(func(frameTime float64) {
		select {
		case fn := <-app.UpQueued:

			fn()
		default:
		}
	})
}
