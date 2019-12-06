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
	// Construct a package context here and use it to sanity-check some things.
	// It also makes a nice cup-holder for the local/remote repositories.
	txCtx := ccmodupdater.PackageTXContext{
		LocalPackages: app.gameInstance.Packages(),
		RemotePackages: middle.GetRemotePackages(),
	}
	// Get local, remote and latest packages for reference.
	localPkg := txCtx.LocalPackages[pkg]
	remotePkg := txCtx.RemotePackages[pkg]
	var latestPkg ccmodupdater.Package = middle.GetLatestOf(localPkg, remotePkg)
	
	// No latest package = no information.
	if latestPkg == nil {
		if middle.InternetConnectionWarning {
			app.MessageBox("Package not available", "The package '" + pkg + "' could not be found.\n\nAs you have ended up here, the package probably had to exist in some form.\nThis error is probably because CCUpdaterUI was unable to retrieve remote packages.\n\n1. Check your internet connection\n2. Try restarting CCUpdaterUI\n3. Contact us", back)
		} else {
			app.MessageBox("I just don't know what went wrong...", "The package '" + pkg + "' could not be found.\nYou should never be able to see this dialog in normal operation.", back)
		}
		return
	}
	
	// Ok, now let's actually start work on the UI
	showInstallButton := false
	annotations := "\n    ID: " + pkg + "\n    Latest Version: " + latestPkg.Metadata().Version().Original()
	if localPkg != nil {
		if remotePkg != nil && remotePkg.Metadata().Version().GreaterThan(localPkg.Metadata().Version()) {
			annotations += "\n    Installed: " + localPkg.Metadata().Version().Original()
			showInstallButton = true
		} else {
			annotations += "\n    Installed"
		}
	} else {
		showInstallButton = true
	}
	chunks := []integration.TypeChunk{
		integration.NewColouredTextTypeChunk(latestPkg.Metadata().HumanName(), design.GlobalFont, design.ThemeText),
		integration.NewColouredTextTypeChunk(annotations, design.ListItemSubTextFont, design.ThemeSubText),
	}
	buttons := []framework.UILayoutElement{}
	if (localPkg != nil) && pkg != "Simplify" {
		removeTx := ccmodupdater.PackageTX{
			pkg: ccmodupdater.PackageTXOperationRemove,
		}
		removeTheme := design.ThemeRemoveActionButton
		_, removeErr := txCtx.Solve(removeTx)
		buttonText := "REMOVE"
		if removeErr != nil {
			buttonText = "NOT REMOVABLE"
			removeTheme = design.ThemeImpossibleActionButton
		}
		buttons = append(buttons, design.ButtonAction(removeTheme, buttonText, func () {
			app.GSDownwards()
			app.PerformTransaction(func () {
				app.GSUpwards()
				app.ShowPackageView(back, pkg)
			}, removeTx)
		}))
	}
	if showInstallButton {
		installTx := ccmodupdater.PackageTX{
			pkg: ccmodupdater.PackageTXOperationInstall,
		}
		buttonText := "INSTALL"
		buttonColour := design.ThemeOkActionButton
		if localPkg != nil {
			buttonText = "UPDATE"
			buttonColour = design.ThemeUpdateActionButton
		}
		_, removeErr := txCtx.Solve(installTx)
		if removeErr != nil {
			buttonText = "NOT INSTALLABLE"
			buttonColour = design.ThemeImpossibleActionButton
		}
		buttons = append(buttons, design.ButtonAction(buttonColour, buttonText, func () {
			app.GSDownwards()
			app.PerformTransaction(func () {
				app.GSUpwards()
				app.ShowPackageView(back, pkg)
			}, installTx)
		}))
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
				Element: design.ButtonBar(buttons),
			},
		},
	})
	
	app.Teleport(design.LayoutDocument(design.Header{
		Title: latestPkg.Metadata().HumanName(),
		Back: back,
	}, fullPanel, true))
}
