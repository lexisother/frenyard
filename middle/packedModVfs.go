package middle

import (
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/ccmod"
	"strings"
)

// PackedModLocation
type PackedModLocation struct {
	// True for packed mods. This overrules everything else, including Drive.
	Valid bool
	Location string
	Metadata ccmodupdater.PackageMetadata
	// If not empty, this is a drive, and has a special overridden name
	Drive string
}

// CheckPackedModLocation converts a path into a PackedModLocation, making sure to set Valid, Location and Metadata correctly in the process.
func CheckPackedModLocation(path string) PackedModLocation {
	// Only try if the location is actually real and looks like a ccmod.
	if BrowserVFSLocationReal(path) && strings.HasSuffix(path, ".ccmod") {
		metadata, err := ccmod.GetMetadata(path)
		if err == nil {
			err = metadata.Verify()
			if err == nil {
				return PackedModLocation{
					Valid: true,
					Location: path,
					Metadata: metadata,
				}
			}
		}
	}
	return PackedModLocation{
		Location: path,
	}
}

// PackedModLocationVFSList performs BrowserVFSList and converts the results to PackedModLocations.
func PackedModFinderVFSList(vfsPath string) []PackedModLocation {
	browserVfs := BrowserVFSList(vfsPath)
	if vfsPath == BrowserVFSPathDefault {
		browserVfs = BrowserVFSAppendDownloads(browserVfs)
	}
	vfsEntries := []PackedModLocation{}
	for _, fi := range browserVfs {
		if fi.Dir {
			vfsEntries = append(vfsEntries, PackedModLocation {
				Location: fi.Location,
				Drive: fi.Drive,
			})
		} else {
			pm := CheckPackedModLocation(fi.Location)
			if pm.Valid {
				vfsEntries = append(vfsEntries, pm)
			}
		}
	}
	return vfsEntries
}

