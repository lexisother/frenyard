package main

import (
	//"fmt"
	"github.com/20kdc/CCUpdaterUI/design"
	//"github.com/20kdc/CCUpdaterUI/frenyard"
	//"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
)

func (app *upApplication) ShowGameFinder() {
	content := design.LayoutMsgbox("It'd be great if this UI actually existed", func () {
		app.ShowGameFinder()
	})
	primary := design.LayoutDocument(design.Header{
		Title: "Enter CrossCode's location",
	}, content, false)
	app.slideContainer.TransitionTo(primary, 1.0, false, false)
	
}
