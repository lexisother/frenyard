package middle

// GameFinderVFSPathDefault is the path for the 'drives panel'. Should be virtual.
const GameFinderVFSPathDefault = "computer://"

// GameFinderVFSEntry reports the analysis results.
type GameFinderVFSEntry struct {
	GameValid bool
	// "Location" of this is always valid. If GameValid is true, the Location must NOT be virtual.
	GameLocation
}

// GameFinderVFSList lists potentially-virtual directories at each path.
func GameFinderVFSList(vfsPath string) []GameFinderVFSEntry {
	if vfsPath == GameFinderVFSPathDefault {
		// Determine drives (OS-dependent)
	}
	return []GameFinderVFSEntry{}
}
