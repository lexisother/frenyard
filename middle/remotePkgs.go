package middle

import (
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

// FakeError should be enabled to prevent internet access by CCUpdaterUI.
const FakeError bool = true

// InternetConnectionWarning is true if the last GetRemotePackages() call actually resulted in error.
var InternetConnectionWarning bool = true

// GetRemotePackages retrieves remote packages from the server. (The CCUpdaterCLI-level cache semantics still apply.)
func GetRemotePackages() map[string]ccmodupdater.RemotePackage {
	InternetConnectionWarning = true
	if !FakeError {
		remote, err := remote.GetRemotePackages()
		if err == nil {
			InternetConnectionWarning = false
			return remote
		}
	}
	return map[string]ccmodupdater.RemotePackage{}
}
