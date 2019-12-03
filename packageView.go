package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

// ShowPackageView shows a dialog for a package.
func (app *upApplication) ShowPackageView(back framework.ButtonBehavior, pkg string) {
	
	localPkg := app.gameInstance.Packages()[pkg]
	remotePkg := middle.GetRemotePackages()[pkg]
	var latestPkg ccmodupdater.Package = localPkg
	if remotePkg != nil {
		if remotePkg.Metadata().Version().GreaterThan(localPkg.Metadata().Version()) {
			latestPkg = remotePkg
		}
	}
	if latestPkg == nil {
		if middle.InternetConnectionWarning {
			app.MessageBox("Package not available", "The package '" + pkg + "' could not be found.\n\nAs you have ended up here, the package probably had to exist in some form.\nThis error is probably because CCUpdaterUI was unable to retrieve remote packages.\n\n1. Check your internet connection\n2. Try restarting CCUpdaterUI\n3. Contact us", back)
		} else {
			app.MessageBox("I just don't know what went wrong...", "The package '" + pkg + "' could not be found.\nYou should never be able to see this dialog in normal operation.", back)
		}
		return
	}
	detail := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				
			},
		},
	})
	app.Teleport(design.LayoutDocument(design.Header{
		Title: latestPkg.Metadata().HumanName(),
		Back: back,
	}, detail, true))
}
