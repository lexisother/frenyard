package main

import (
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
)

// The Primary View
func (app *upApplication) ShowPrimaryView() {
	app.ShowWaiter(framework.SlideTransition{
		Length: 1.0,
	}, "NYI", func (f func(string)) {

	}, func () {
		
	});
}
