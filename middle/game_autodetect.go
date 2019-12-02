package middle

import (
	"os"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
)

// GameLocation represents a location a copy of CrossCode is in.
type GameLocation struct {
	Location string
	Version string
}

// detectPossibleGameLocations returns a list of likely game locations.
func detectPossibleGameLocations() []string {
	home := os.Getenv("HOME")
	locations := []string{
		home + "/.steam/steam/steamapps/common/CrossCode",
		// Sorry, but I need to test stacking these
		home + "/Documents/Local/Quicksand/apple",
	}
	// Windows users may have a home drive from A to Z.
	// A and B are historically floppy disk drives, but try anyway.
	for i := 0; i < 26; i++ {
		drive := string(rune('A' + i))
		locations = append(locations, drive + ":\\Program Files/Steam/steamapps/common/CrossCode", drive + ":\\Program Files (x86)/Steam/steamapps/common/CrossCode")
	}
	return locations
}

// AutodetectGameLocations attempts to find game locations. This may take some time as it can spin the CD/DVD drive on Windows; run in a Goroutine!
func AutodetectGameLocations() []GameLocation {
	possibilities := detectPossibleGameLocations()
	locations := []GameLocation{}
	for _, v := range possibilities {
		gameInstance := ccmodupdater.NewGameInstance(v)
		plugins, err := local.AllLocalPackagePlugins(gameInstance)
		if err != nil {
			continue
		}
		gameInstance.LocalPlugins = plugins
		ccDetails, hasCC := gameInstance.Packages()["crosscode"]
		if hasCC {
			locations = append(locations, GameLocation{
				v,
				ccDetails.Metadata().Version().Original(),
			})
		}
	}
	return locations
}
