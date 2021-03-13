package middle

import (
	"os"
	"io/ioutil"
	"github.com/20kdc/go-vkv"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
	"unicode/utf8"
	"strconv"
	"fmt"
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

// detectPossibleSteamRoots returns a list of Steam libraries, not including additional manually added libraries.
func detectPossibleSteamRoots() []string {
	home := os.Getenv("HOME")
	locations := []string{
		home + "/.steam/steam",
		// Thanks to dmitmel for this: Mac OS X location
		home + "/Library/Application Support/Steam",
	}
	// Windows users may have a home drive from A to Z.
	// A and B are historically floppy disk drives, but try anyway.
	for i := 0; i < 26; i++ {
		drive := string(rune('A' + i))
		locations = append(locations, drive + ":\\Program Files/Steam", drive + ":\\Program Files (x86)/Steam")
	}
	return locations
}

// appendHarvestFromSteamLibraryFolders is the inner critical part of detectSteamLeaves.
func appendHarvestFromSteamLibraryFolders(locations []string, data []byte, where string) []string {
	if !utf8.Valid(data) {
		// Not strictly sure how Go would take this; avoid having to answer the question for now
		fmt.Printf(" invalid: not UTF-8\n")
		return locations
	}
	tkns, err := kvkv.InTokenize(string(data), true, where)
	if err != nil {
		fmt.Printf(" invalid (tkn): %s\n", err)
		return locations
	}
	obj, err := kvkv.InParse(tkns)
	if err != nil {
		fmt.Printf(" invalid (par): %s\n", err)
		return locations
	}
	// we may have something
	subObj, err := obj.FindObject("LibraryFolders")
	if err != nil {
		fmt.Printf(" invalid (lff): %s\n", err)
		return locations
	}
	index := 1
	for {
		lib, err := subObj.FindString(strconv.Itoa(index))
		if err != nil {
			return locations
		}
		fmt.Printf(" found: %s\n", lib)
		locations = append(locations, lib)
		index++
	}
}

// detectSteamLeaves expands a detectPossibleSteamRoots list to include alternate library locations.
func detectSteamLeaves() []string {
	locations := detectPossibleSteamRoots()
	for _, v := range locations {
		path := v + "/steamapps/libraryfolders.vdf"
		file, err := os.Open(path)
		if err == nil {
			fmt.Printf("Steam root: " + v + " is valid & has library folders index...\n")
			bytes, err := ioutil.ReadAll(file)
			if err == nil {
				locations = appendHarvestFromSteamLibraryFolders(locations, bytes, path)
			}
			file.Close()
		}
	}
	return locations
}

// detectPossibleGameLocations returns a list of likely game locations.
func detectPossibleGameLocations() []string {
	locations := []string{}
	for _, v := range detectSteamLeaves() {
		// Thanks to dmitmel for the Mac OS X location
		locations = append(locations, v + "/steamapps/common/CrossCode", v + "/steamapps/common/CrossCode/CrossCode.app/Contents/Resources/app.nw")
	}
	return locations
}

// CheckGameLocation converts a path into a GameLocation, making sure to set Valid & Version correctly in the process.
func CheckGameLocation(dir string) GameLocation {
	// Only try if the location is actually real.
	if BrowserVFSLocationReal(dir) {
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
