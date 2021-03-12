package middle

import (
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
	"github.com/20kdc/CCUpdaterUI/design"
)

// FakeError should be enabled to prevent internet access by CCUpdaterUI.
const FakeError bool = false

// InternetConnectionWarning is true if the last GetRemotePackages() call actually resulted in error.
var InternetConnectionWarning bool = true

// GetRemotePackages retrieves remote packages from the server. (The CCUpdaterCLI-level cache semantics still apply.)
func GetRemotePackages() map[string]ccmodupdater.RemotePackage {
	InternetConnectionWarning = true
	if !FakeError {
		updateAlertHook()
		remote, err := remote.GetRemotePackages()
		if err == nil {
			InternetConnectionWarning = false
			return remote
		}
	}
	return map[string]ccmodupdater.RemotePackage{}
}

// GetLatestOf returns the latest of two possibly-nil packages (returning nil if both are nil)
func GetLatestOf(local ccmodupdater.Package, remote ccmodupdater.Package) ccmodupdater.Package {
	if local != nil {
		if remote != nil {
			if remote.Metadata().Version().GreaterThan(local.Metadata().Version()) {
				return remote
			}
		}
		return local
	}
	return remote
}

// PackageIcon returns the relevant icon ID for a package.
func PackageIcon(pkg ccmodupdater.Package) design.IconID {
	typ := pkg.Metadata().Type()
	if typ == ccmodupdater.PackageTypeMod {
		return design.ModIconID
	} else if typ == ccmodupdater.PackageTypeTool {
		return design.ToolIconID
	}
	return design.DirectoryIconID
}

// PackagePVIcon returns the relevant icon ID for the Primary View given the state of the package locally and remotely.
func PackagePVIcon(local ccmodupdater.Package, remote ccmodupdater.Package) design.IconID {
	if local != nil {
		if remote != nil {
			if remote.Metadata().Version().GreaterThan(local.Metadata().Version()) {
				return design.UpdatableIconID
			}
		}
		typ := local.Metadata().Type()
		if typ == ccmodupdater.PackageTypeMod {
			return design.InstalledIconID
		}
		return PackageIcon(local)
	}
	return design.BlankIconID
}
