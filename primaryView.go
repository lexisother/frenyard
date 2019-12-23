package main

import (
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"sort"
	"os"
)

// ShowPrimaryView shows the "Primary View" (the mod list right now)
func (app *upApplication) ShowPrimaryView() {
	
	if app.cachedPrimaryView != nil {
		app.Teleport(app.cachedPrimaryView)
		return
	}
	
	slots := []framework.FlexboxSlot{}
	
	warnings := middle.FindWarnings(app.gameInstance)
	if app.config.DevMode {
		warnings = append(warnings, middle.Warning{
			Text: "You are in developer mode! Go to the Build Information (top-right button, 'Credits', 'Build Information') to deactivate it.",
			Action: middle.NullActionWarningID,
		})
	}
	for _, v := range warnings {
		fixAction := framework.ButtonBehavior(nil)
		if v.Action == middle.InstallOrUpdatePackageWarningID {
			pkgID := v.Parameter
			fixAction = func () {
				app.GSRightwards()
				app.ShowPackageView(func () {
					app.GSLeftwards()
					app.ShowPrimaryView()
				}, pkgID)
			}
		} else if v.Action == middle.URLAndCloseWarningID {
			url := v.Parameter
			fixAction = func () {
				middle.OpenURL(url)
				os.Exit(0)
			}
		}
		slots = append(slots, framework.FlexboxSlot{
			Element: design.InformationPanel(design.InformationPanelDetails{
				Text: v.Text,
				ActionText: "FIX",
				Action: fixAction,
			}),
		})
	}
	
	// Ok, let's get all the packages in a nice row
	localPackages := app.gameInstance.Packages()
	remotePackages := middle.GetRemotePackages()
	packageSet := make(map[string]bool)
	packageList := []design.ListItemDetails{}
	for k := range localPackages {
		packageSet[k] = true
	}
	for k := range remotePackages {
		packageSet[k] = true
	}
	// Actually build the UI now!
	for pkgID := range packageSet {
		local := localPackages[pkgID]
		remote := remotePackages[pkgID]
		latest := middle.GetLatestOf(local, remote)
		var typeCheck ccmodupdater.Package = local
		if typeCheck == nil {
			typeCheck = remote
		}
		
		if (!app.config.DevMode) && (typeCheck.Metadata().Type() != ccmodupdater.PackageTypeMod) {
			continue
		}
		
		status := "unable to comprehend status"
		if local != nil && remote != nil {
			lmv := local.Metadata().Version()
			rmv := remote.Metadata().Version()
			if local.Metadata().Version().GreaterThan(remote.Metadata().Version()) {
				status = lmv.Original() + " installed (local development build, " + rmv.Original() + " remote)"
			} else if remote.Metadata().Version().GreaterThan(local.Metadata().Version()) {
				status = lmv.Original() + " installed (out of date, " + rmv.Original() + " available)"
			} else {
				status = lmv.Original() + " (up to date)"
			}
		} else if local != nil {
			status = latest.Metadata().Version().Original() + " installed (no remote copy)"
		} else if remote != nil {
			status = latest.Metadata().Version().Original() + " available"
		}
		pkgIDLocal := pkgID
		packageList = append(packageList, design.ListItemDetails{
			Icon: middle.PackageIcon(latest),
			Text: latest.Metadata().HumanName(),
			Subtext: status,
			Click: func () {
				app.GSRightwards()
				app.ShowPackageView(func () {
					app.GSLeftwards()
					app.ShowPrimaryView()
				}, pkgIDLocal)
			},
		})
	}

	sort.Sort(design.SortListItemDetails(packageList))
	slots = append(slots, framework.FlexboxSlot{
		Element: design.NewUISearchBoxPtr("Search...", packageList),
		Grow: 1,
	})

	// This cached element is used to boost performance when possible.
	app.cachedPrimaryView = design.LayoutDocument(design.Header{
		Title: "Mods",
		Back: func () {
			app.cachedPrimaryView = nil
			app.GSLeftwards()
			app.ResetWithGameLocation(false, middle.GameFinderVFSPathDefault)
		},
		BackIcon: design.GameIconID,
		ForwardIcon: design.MenuIconID,
		Forward: func () {
			app.GSRightwards()
			app.ShowOptionsMenu(func () {
				app.GSLeftwards()
				app.ShowPrimaryView()
			})
		},
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: slots,
	}), true)
	app.Teleport(app.cachedPrimaryView)
}
