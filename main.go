package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

type upApplication struct {
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
