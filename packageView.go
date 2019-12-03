package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
)

// ShowPackageView shows a dialog for a package.
func (app *upApplication) ShowPackageView(back framework.ButtonBehavior, pkg string) {
	
	localPkg := app.gameInstance.Packages()[pkg]
	remotePkg := middle.GetRemotePackages()[pkg]
	var latestPkg ccmodupdater.Package = middle.GetLatestOf(localPkg, remotePkg)
	if latestPkg == nil {
		if middle.InternetConnectionWarning {
			app.MessageBox("Package not available", "The package '" + pkg + "' could not be found.\n\nAs you have ended up here, the package probably had to exist in some form.\nThis error is probably because CCUpdaterUI was unable to retrieve remote packages.\n\n1. Check your internet connection\n2. Try restarting CCUpdaterUI\n3. Contact us", back)
		} else {
			app.MessageBox("I just don't know what went wrong...", "The package '" + pkg + "' could not be found.\nYou should never be able to see this dialog in normal operation.", back)
		}
		return
	}
	
	annotations := "\n    ID: " + pkg + "\n    Version: " + latestPkg.Metadata().Version().Original()
	if localPkg != nil {
		if latestPkg == remotePkg {
			annotations += "\n    Installed (Outdated)"
		} else {
			annotations += "\n    Installed"
		}
	}
	chunks := []integration.TypeChunk{
		integration.NewColouredTextTypeChunk(latestPkg.Metadata().HumanName(), design.GlobalFont, design.ThemeText),
		integration.NewColouredTextTypeChunk(annotations, design.ListItemSubTextFont, design.ThemeSubText),
	}
	
	detail := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: false,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: design.NewIconPtr(0xFFFFFFFF, middle.PackageIcon(latestPkg), 36),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewCompoundTypeChunk(chunks), 0xFFFFFFFF, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignStart}),
			},
		},
	})
	
	fullPanel := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: detail,
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk(latestPkg.Metadata().Description(), design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{X: frenyard.AlignStart, Y: frenyard.AlignStart}),
				Shrink: 1,
			},
			framework.FlexboxSlot{
				Grow: 1,
				Shrink: 1,
			},
			framework.FlexboxSlot{
			},
		},
	})
	
	app.Teleport(design.LayoutDocument(design.Header{
		Title: latestPkg.Metadata().HumanName(),
		Back: back,
	}, fullPanel, true))
}
