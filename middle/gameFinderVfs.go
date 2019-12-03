package middle

import (
	"os"
	"runtime"
	"io/ioutil"
	"path/filepath"
)

// GameFinderVFSPathDefault is the path for the 'drives panel'. Should be virtual.
const GameFinderVFSPathDefault = "computer://"

// GameFinderVFSList lists potentially-virtual directories at each path.
func GameFinderVFSList(vfsPath string) []GameLocation {
	vfsEntries := []GameLocation{}

	if vfsPath == GameFinderVFSPathDefault {
		// Determine drives (OS-dependent)
		if runtime.GOOS == "windows" {
			// we can get away with this, right?
			for i := 0; i < 26; i++ {
				drive := string(rune('A' + i)) + ":\\"
				_, err := ioutil.ReadDir(drive)
				if err == nil {
					vfsEntries = append(vfsEntries, GameLocation{
						Drive: drive,
						Location: drive,
					})
				}
			}
			return vfsEntries
		}
		home := os.Getenv("HOME")
		return []GameLocation{
			GameLocation{
				Drive: "Home",
				Location: home,
			},
			GameLocation{
				Drive: "Root",
				Location: "/",
			},
			GameLocation{
				Drive: "PWD",
				Location: ".",
			},
		}
	}
	
	fileInfos, err := ioutil.ReadDir(vfsPath)
	if err == nil {
		for _, fi := range fileInfos {
			if fi.IsDir() {
				vfsEntries = append(vfsEntries, CheckGameLocation(filepath.Join(vfsPath, fi.Name())))
			}
		}
	}
	
	// List
	return vfsEntries
}
