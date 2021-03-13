package middle

// GameFinderVFSList performs BrowserVFSList and converts the results to GameLocations.
func GameFinderVFSList(vfsPath string) []GameLocation {
	vfsEntries := []GameLocation{}
	for _, fi := range BrowserVFSList(vfsPath) {
		if fi.Dir {
			cgl := CheckGameLocation(fi.Location)
			cgl.Drive = fi.Drive
			vfsEntries = append(vfsEntries, cgl)
		}
	}
	return vfsEntries
}
