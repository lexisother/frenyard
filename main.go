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
	slideContainer *framework.UISlideTransitionContainer
	window frenyard.Window
	upQueued chan func()
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
		slideContainer: slideContainer,
		window: wnd,
		upQueued: make(chan func(), 16),
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
