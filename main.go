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
	// Reset this to nil whenever anything that affects the Primary View changes
	cachedPrimaryView framework.UILayoutElement
	config middle.UpdaterConfig
	mainContainer *framework.UISlideTransitionContainer
	window frenyard.Window
	upQueued chan func()
	teleportSettings framework.SlideTransition
}

const upTeleportLen float64 = 0.25

// GSLeftwards sets the teleportation affinity to LEFT.
func (app *upApplication) GSLeftwards() {
	app.teleportSettings.Reverse = true
	app.teleportSettings.Vertical = false
	app.teleportSettings.Length = upTeleportLen
}
// GSRightwards sets the teleportation affinity to RIGHT.
func (app *upApplication) GSRightwards() {
	app.teleportSettings.Reverse = false
	app.teleportSettings.Vertical = false
	app.teleportSettings.Length = upTeleportLen
}
// GSLeftwards sets the teleportation affinity to UP.
func (app *upApplication) GSUpwards() {
	app.teleportSettings.Reverse = true
	app.teleportSettings.Vertical = true
	app.teleportSettings.Length = upTeleportLen
}
// GSRightwards sets the teleportation affinity to DOWN.
func (app *upApplication) GSDownwards() {
	app.teleportSettings.Reverse = false
	app.teleportSettings.Vertical = true
	app.teleportSettings.Length = upTeleportLen
}
// GSInstant sets the teleportation affinity to INSTANT.
func (app *upApplication) GSInstant() {
	// direction doesn't matter
	app.teleportSettings.Length = 0
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
