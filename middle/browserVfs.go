package middle

import (
	"os"
	"runtime"
	"io/ioutil"
	"path/filepath"
)

// BrowserLocation represents a location in the browser VFS.
type BrowserLocation struct {
	// Location is the actual path (should be passable to file IO functions if not virtual)
	Location string
	// Dir is true if this is a directory, false otherwise.
	Dir bool
	// If not empty, this is a drive, and has a special overridden name
	Drive string
}

// BrowserVFSPathDefault is the path for the 'drives panel'. Should be virtual.
const BrowserVFSPathDefault = "computer://"

// BrowserVFSLocationReal returns true if the path is a real file or directory that can be accessed using standard file IO.
func BrowserVFSLocationReal(vfsPath string) bool {
	return vfsPath != BrowserVFSPathDefault
}

// BrowserVFSList lists potentially-virtual directories at each path.
func BrowserVFSList(vfsPath string) []BrowserLocation {
	vfsEntries := []BrowserLocation{}

	if vfsPath == BrowserVFSPathDefault {
		// Determine drives (OS-dependent)
		if runtime.GOOS == "windows" {
			// we can get away with this, right?
			for i := 0; i < 26; i++ {
				drive := string(rune('A' + i)) + ":\\"
				_, err := ioutil.ReadDir(drive)
				if err == nil {
					vfsEntries = append(vfsEntries, BrowserLocation {
						Drive: drive,
						Location: drive,
						Dir: true,
					})
				}
			}
			return vfsEntries
		}
		home := os.Getenv("HOME")
		return []BrowserLocation {
			BrowserLocation {
				Drive: "Home",
				Location: home,
				Dir: true,
			},
			BrowserLocation {
				Drive: "Root",
				Location: "/",
				Dir: true,
			},
			BrowserLocation {
				Drive: "PWD",
				Location: ".",
				Dir: true,
			},
		}
	}
	
	fileInfos, err := ioutil.ReadDir(vfsPath)
	if err == nil {
		for _, fi := range fileInfos {
			vfsEntries = append(vfsEntries, BrowserLocation {
				Location: filepath.Join(vfsPath, fi.Name()),
				Dir: fi.IsDir(),
			})
		}
	}
	
	// List
	return vfsEntries
}

// Appends the Downloads directory
func BrowserVFSAppendDownloads(existing []BrowserLocation) []BrowserLocation {
	home := os.Getenv("HOME")
	uprof := os.Getenv("USERPROFILE")
	locations := []string{
		home + "/Downloads",
		uprof + "/Downloads",
	}
	for _, v := range locations {
		_, err := ioutil.ReadDir(v)
		if err == nil {
			existing = append(existing, BrowserLocation {
				Drive: "Downloads",
				Location: v,
				Dir: true,
			})
			break
		}
	}
	return existing
}
