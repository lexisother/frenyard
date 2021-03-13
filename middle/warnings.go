package middle

import (
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/Masterminds/semver"
	"net/http"
	"io/ioutil"
)

// autoupdate alert stuff ; update this for new version names!!!
// version schedule:
// past: "lea\n" "emilie\n"
// present: "" (disable)
// future: ???
// note that the "\n" is important if this is a real ID
const localUIVersionID string = ""
var remoteUIVersionID string = localUIVersionID
var hasAlreadyCheckedRemoteUIVersionID bool = false

func updateAlertHook() {
	if localUIVersionID == "" {
		// version checks are disabled
		return
	}
	if !hasAlreadyCheckedRemoteUIVersionID {
		hasAlreadyCheckedRemoteUIVersionID = true
	} else {
		return
	}
	res, err := http.Get("https://20kdc.duckdns.org/ccmodloader/version")
	if err == nil {
		data, err := ioutil.ReadAll(res.Body)
		if err == nil {
			remoteUIVersionID = string(data)
		}
	}
}

// WarningID represents a kind of warning action.
type WarningID int

const (
	// NullActionWarningID cannot be automatically fixed.
	NullActionWarningID WarningID = iota
	// InstallOrUpdatePackageWarningID warnings can be solved by installing/updating the package Parameter.
	InstallOrUpdatePackageWarningID
	// URLAndCloseWarningID warnings can be solved manually by the user given navigation to a URL. The application closes.
	URLAndCloseWarningID
)

// Warning represents a warning to show the user on the primary view.
type Warning struct {
	Text string
	Action WarningID
	Parameter string
}

// FindWarnings detects issues with the installation to show on the primary view.
func FindWarnings(game *ccmodupdater.GameInstance) []Warning {
	warnings := []Warning{}
	if InternetConnectionWarning {
		warnings = append(warnings, Warning{
			Text: "CCUpdaterUI wasn't able to retrieve the mod metadata; downloading mods is not possible.",
		})
	}
	if localUIVersionID != remoteUIVersionID {
		warnings = append(warnings, Warning{
			Text: "CCUpdaterUI is out of date! Please update it.",
			Action: URLAndCloseWarningID,
			Parameter: "https://20kdc.duckdns.org/ccmodloader/update-thunk.html",
		})
	}
	crosscode := game.Packages()["crosscode"]
	if crosscode == nil {
		warnings = append(warnings, Warning{
			Text: "CrossCode is not installed in this installation. (Ok, come on, how'd you manage this? - 20kdc)",
		})
	} else if crosscode.Metadata().Version().LessThan(semver.MustParse("1.1.0")) {
		warnings = append(warnings, Warning{
			Text: "The CrossCode version is " + crosscode.Metadata().Version().Original() + "; mods usually expect 1.1.0 or higher.",
		})
	} 
	ccloader := game.Packages()["ccloader"]
	if ccloader == nil {
		warnings = append(warnings, Warning{
			Text: "No modloader is installed; thus any mods installed cannot be run.",
			Action: InstallOrUpdatePackageWarningID,
			Parameter: "ccloader",
		})
	} else {
		remoteCCLoader := GetRemotePackages()["ccloader"]
		if remoteCCLoader != nil {
			if remoteCCLoader.Metadata().Version().GreaterThan(ccloader.Metadata().Version()) {
				warnings = append(warnings, Warning{
					Text: "CCLoader is out of date. This may cause buggy behavior, or mods may rely on missing features.",
					Action: InstallOrUpdatePackageWarningID,
					Parameter: "ccloader",
				})
			}
		}
	}
	return warnings
}
