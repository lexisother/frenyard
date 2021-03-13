package main

import (
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/20kdc/CCUpdaterUI/frenyard"
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/frenyard/integration"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

type upPackedModState map[string]ccmodupdater.RemotePackage

func (app *upApplication) ShowInitPackedModManager(back framework.ButtonBehavior, firstModPath string) {
	initState := make(upPackedModState)
	app.ShowPackedModManagerAttemptToAddMod(initState, back, func () {
		app.ShowPackedModManager(back, initState)
	}, firstModPath)
}

func (app *upApplication) ShowPackedModManagerAttemptToAddMod(state upPackedModState, back framework.ButtonBehavior, fwd framework.ButtonBehavior, modPath string) {
	pmr, err := remote.NewPackedModRemotePackage(modPath)
	if err != nil {
		app.MessageBox("Error", err.Error(), back)
	}
	id := pmr.Metadata().Name()
	if state[id] != nil {
		app.MessageBox("Error", "You've already added mod " + id + ".", back)
	} else {
		state[id] = pmr
		fwd()
	}
}

func (app *upApplication) ShowPackedModManager(back framework.ButtonBehavior, mods upPackedModState) {
	backHere := func () {
		app.GSLeftwards()
		app.ShowPackedModManager(back, mods)
	}

	localPackages := app.gameInstance.Packages()
	remotePackages := make(map[string]ccmodupdater.RemotePackage)
	// Overlay these two to create the actual transaction environment
	for k, v := range middle.GetRemotePackages() {
		remotePackages[k] = v
	}
	for k, v := range mods {
		remotePackages[k] = v
	}

	packageTXContext := ccmodupdater.PackageTXContext {
		LocalPackages: localPackages,
		RemotePackages: remotePackages,
	}

	// Setup the transaction, solve it, and find involved of results
	packageTX := make(ccmodupdater.PackageTX)
	for modID, _ := range mods {
		packageTX[modID] = ccmodupdater.PackageTXOperationInstall
	}
	// This relies on packageTXSolutions always containing some meaningful information even on error.
	packageTXSolutions, packageTXErr := packageTXContext.Solve(packageTX)
	involved := make(map[string]bool)
	for _, solution := range packageTXSolutions {
		for k, _ := range solution {
			involved[k] = true
		}
	}

	modSlots := []framework.FlexboxSlot{}
	for modID, modPkg := range mods {
		localModID := modID
		modSlots = append(modSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Icon: design.ModIconID,
				Text: modID,
				Subtext: modPkg.Metadata().Description(),
				Click: func () {
					delete(mods, localModID)
					backHere()
				},
			}),
			RespectMinimumSize: true,
		})
	}
	for depID, _ := range involved {
		if mods[depID] != nil {
			// ok, it's already dealt with
			continue
		}
		// it hasn't been dealt with, show an obvious warning
		depItem := design.ListItemDetails{
			Icon: design.WarningIconID,
			Text: depID,
			Subtext: "This is missing. If it's available as a .ccmod, you can add it - it isn't available in CCModDB.",
			Click: func () {
				// Nothing to do here, can't do anything...
			},
		}
		if remotePackages[depID] != nil {
			depItem.Icon = design.UpdatableIconID
			depItem.Subtext = "This will automatically be downloaded. Alternatively, you can add a specific version as a .ccmod here."
		}
		modSlots = append(modSlots, framework.FlexboxSlot{
			Element: design.ListItem(depItem),
			RespectMinimumSize: true,
		})
	}
	// Space-taker to prevent wrongly scaled list items
	modSlots = append(modSlots, framework.FlexboxSlot{
		Grow: 1,
		Shrink: 0,
	})

	readyModsScroller := design.ScrollboxV(framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		WrapMode: framework.FlexboxWrapModeNone,
		Slots: modSlots,
	}))

	installButtonColour := design.ThemeOkActionButton
	installButtonText := "INSTALL"

	if packageTXErr != nil {
		installButtonColour = design.ThemeImpossibleActionButton
		installButtonText = "NOT INSTALLABLE"
	}

	content := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			framework.FlexboxSlot{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("This is the set of mods to install so far. To remove a mod, click on it in the list.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: readyModsScroller,
				Grow: 1,
				Shrink: 1,
				RespectMinimumSize: true,
			},
			framework.FlexboxSlot{
				Basis: design.SizeMarginAroundEverything,
			},
			framework.FlexboxSlot{
				Element: design.ButtonBar([]framework.UILayoutElement{
					design.ButtonAction(design.ThemeUpdateActionButton, "ADD MOD...", func () {
						app.GSRightwards()
						app.ShowPackedModFinder(backHere, func (path string) {
							app.ShowPackedModManagerAttemptToAddMod(mods, backHere, backHere, path)
						}, middle.BrowserVFSPathDefault)
					}),
					design.ButtonAction(installButtonColour, installButtonText, func () {
						app.GSRightwards()
						app.PerformTransaction(func (success bool) {
							if success {
								app.GSLeftwards()
								app.ShowPrimaryView()
							} else {
								backHere()
							}
						}, packageTX, remotePackages)
					}),
				}),
			},
		},
	})
	primary := design.LayoutDocument(design.Header{
		Back: back,
		Title: "Install Packed Mods...",
	}, content, true)
	app.Teleport(primary)
}
