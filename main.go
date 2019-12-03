package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

type upApplication struct {
	gameInstance *ccmodupdater.GameInstance
	config middle.UpdaterConfig
	mainContainer *framework.UISlideTransitionContainer
	window frenyard.Window
	upQueued chan func()
	teleportSettings framework.SlideTransition
}

// GSLeftwards sets the teleportation affinity to LEFT.
func (app *upApplication) GSLeftwards() {
	app.teleportSettings.Reverse = true
	app.teleportSettings.Vertical = false
}
// GSRightwards sets the teleportation affinity to RIGHT.
func (app *upApplication) GSRightwards() {
	app.teleportSettings.Reverse = false
	app.teleportSettings.Vertical = false
}
// GSLeftwards sets the teleportation affinity to UP.
func (app *upApplication) GSUpwards() {
	app.teleportSettings.Reverse = true
	app.teleportSettings.Vertical = true
}
// GSRightwards sets the teleportation affinity to DOWN.
func (app *upApplication) GSDownwards() {
	app.teleportSettings.Reverse = false
	app.teleportSettings.Vertical = true
}
// Teleport starts a transition with the cached affinity settings.
func (app *upApplication) Teleport(target framework.UILayoutElement) {
	forkTD := app.teleportSettings
	forkTD.Element = target
	app.mainContainer.TransitionTo(forkTD)
}

func main() {
	frenyard.TargetFrameTime = 0.016
	slideContainer := framework.NewUISlideTransitionContainerPtr(nil)
	slideContainer.FyEResize(design.SizeWindowInit)
	wnd, err := framework.CreateBoundWindow("CCUpdaterUI", true, design.ThemeBackground, slideContainer)
	if err != nil {
		panic(err)
	}
	design.Setup(frenyard.InferScale(wnd))
	wnd.SetSize(design.SizeWindow)
	// Ok, now get it ready.
	app := (&upApplication{
		config: middle.ReadUpdaterConfig(),
		mainContainer: slideContainer,
		window: wnd,
		upQueued: make(chan func(), 16),
		teleportSettings: framework.SlideTransition{
			Length: 1.0,
		},
	})
	app.ShowGameFinderPreface()
	// Started!
	frenyard.GlobalBackend.Run(func(frameTime float64) {
		select {
			case fn := <- app.upQueued:
				fn()
			default:
		}
	})
}
