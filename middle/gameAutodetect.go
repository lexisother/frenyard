package middle

import (
	"os"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
)

// GameLocation represents a location a copy of CrossCode may be in.
type GameLocation struct {
	// True for CrossCode installs. This overrules everything else, including Drive.
	Valid bool
	Location string
	Version string
	// If not empty, this is a drive, and has a special overridden name
	Drive string
}

// detectPossibleGameLocations returns a list of likely game locations.
func detectPossibleGameLocations() []string {
	home := os.Getenv("HOME")
	locations := []string{
		home + "/.steam/steam/steamapps/common/CrossCode",
		// Thanks to dmitmel for this: Mac OS X location
		home + "/Library/Application Support/Steam/steamapps/common/CrossCode/CrossCode.app/Contents/Resources/app.nw",
	}
	// Windows users may have a home drive from A to Z.
	// A and B are historically floppy disk drives, but try anyway.
	for i := 0; i < 26; i++ {
		drive := string(rune('A' + i))
		locations = append(locations, drive + ":\\Program Files/Steam/steamapps/common/CrossCode", drive + ":\\Program Files (x86)/Steam/steamapps/common/CrossCode")
	}
	return locations
}

// CheckGameLocation converts a path into a GameLocation, making sure to set Valid & Version correctly in the process.
func CheckGameLocation(dir string) GameLocation {
	gameInstance := ccmodupdater.NewGameInstance(dir)
	plugins, err := local.AllLocalPackagePlugins(gameInstance)
	if err == nil {
		gameInstance.LocalPlugins = plugins
		ccDetails, hasCC := gameInstance.Packages()["crosscode"]
		if hasCC {
			return GameLocation{
				Valid: true,
				Location: dir,
				Version: ccDetails.Metadata().Version().Original(),
			}
		}
	}
	return GameLocation{
		Location: dir,
	}
}

// AutodetectGameLocations attempts to find game locations. This may take some time as it can spin the CD/DVD drive on Windows; run in a Goroutine!
func AutodetectGameLocations() []GameLocation {
	possibilities := detectPossibleGameLocations()
	locations := []GameLocation{}
	for _, v := range possibilities {
		loc := CheckGameLocation(v)
		if loc.Valid {
			locations = append(locations, loc)
		}
	}
	return locations
}
