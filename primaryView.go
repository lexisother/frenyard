package main

import (
	"github.com/20kdc/CCUpdaterUI/frenyard/framework"
	"github.com/20kdc/CCUpdaterUI/design"
	"github.com/20kdc/CCUpdaterUI/middle"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"sort"
	"strings"
	"os"
)

type primaryViewPackageListSortable struct {
	names []string
	localPackages map[string]ccmodupdater.LocalPackage
	remotePackages map[string]ccmodupdater.RemotePackage
	packageLocalLatest map[string]bool
	humanNamesLowered map[string]string
}

func (slid primaryViewPackageListSortable) Len() int {
	return len(slid.names)
}
func (slid primaryViewPackageListSortable) Swap(i int, j int) {
	a := slid.names[i]
	slid.names[i] = slid.names[j]
	slid.names[j] = a
}
func (slid primaryViewPackageListSortable) Less(i int, j int) bool {
	iPkgName := slid.names[i]
	jPkgName := slid.names[j]
	// check 1: packages we have locally always go first
	iLoPkg := slid.localPackages[iPkgName]
	jLoPkg := slid.localPackages[jPkgName]
	if iLoPkg != nil && jLoPkg == nil {
		return true
	} else if jLoPkg != nil && iLoPkg == nil {
		return false
	}
	// check 2: packages where local is latest go last
	iLatest := slid.packageLocalLatest[iPkgName]
	jLatest := slid.packageLocalLatest[jPkgName]
	if iLatest && !jLatest {
		return false
	} else if jLatest && !iLatest {
		return true
	}
	res := strings.Compare(slid.humanNamesLowered[iPkgName], slid.humanNamesLowered[jPkgName])
	return res < 0
}

// getPrimaryViewPackageListSortable gets all the data about packages to show, sorts it, etc.
func (app *upApplication) getPrimaryViewPackageListSortable() primaryViewPackageListSortable {
	// The actual input
	localPackages := app.gameInstance.Packages()
	remotePackages := middle.GetRemotePackages()
	// Set (intermediate step used so that if a package is common between local and remote, it gets iterated once)
	packageSet := make(map[string]bool)
	for k := range localPackages {
		packageSet[k] = true
	}
	for k := range remotePackages {
		packageSet[k] = true
	}
	// Expand into what's necessary for the sortable
	names := []string{}
	packageLocalLatest := make(map[string]bool)
	humanNamesLowered := make(map[string]string)
	for k, _ := range packageSet {
		// in
		local := localPackages[k]
		remote := remotePackages[k]
		// work out stuff
		latest := middle.GetLatestOf(local, remote)
		localLatest := false
		if local != nil {
			localLatest = !latest.Metadata().Version().GreaterThan(local.Metadata().Version())
		}
		humanNameLowered := strings.ToLower(latest.Metadata().HumanName())
		// out
		names = append(names, k)
		packageLocalLatest[k] = localLatest
		humanNamesLowered[k] = humanNameLowered
	}
	// The actual struct
	sortable := primaryViewPackageListSortable {
		names: names,
		localPackages: localPackages,
		remotePackages: remotePackages,
		packageLocalLatest: packageLocalLatest,
		humanNamesLowered: humanNamesLowered,
	}
	// Finally, sort & return
	sort.Sort(sortable)
	return sortable
}

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
	packageList := app.getPrimaryViewPackageListSortable()
	packageListItems := []design.ListItemDetails{}
	// Actually build the UI now!
	for _, pkgID := range packageList.names {
		local := packageList.localPackages[pkgID]
		remote := packageList.remotePackages[pkgID]
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
		description := latest.Metadata().Description()
		if description != "" {
			status = description + "\n" + status
		}
		pkgIDLocal := pkgID
		packageListItems = append(packageListItems, design.ListItemDetails{
			Icon: middle.PackagePVIcon(local, remote),
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

	slots = append(slots, framework.FlexboxSlot{
		Element: design.NewUISearchBoxPtr("Search...", packageListItems),
		Grow: 1,
	})

	// This cached element is used to boost performance when possible.
	app.cachedPrimaryView = design.LayoutDocument(design.Header{
		Title: "Mods",
		Back: func () {
			app.cachedPrimaryView = nil
			app.GSLeftwards()
			// The idea here is that BrowserVFSPathDefault is never a valid path.
			app.ResetWithGameLocation(false, middle.BrowserVFSPathDefault)
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
