package middle

import (
	"os"
	"fmt"
	"path/filepath"
)

// Launch launches the game given a gameinstance base directory.
func Launch(base string) (*os.Process, error) {
	// Try to run the game...
	executables := []string{
		// Unixes
		"run", // Override script
		"nw", // NW.JS SDK replacement
		"CrossCode", // Original executable
		"../../MacOS/nwjs", // Mac OS X executable
		// Windows
		"nw.exe", // NW.JS SDK replacement
		"CrossCode.exe", // Original executable
	}
	for _, executable := range executables {
		fullPath := filepath.Join(base, executable)
		proc, err := os.StartProcess(fullPath, []string{fullPath}, &os.ProcAttr{
			Dir: base,
		})
		if err == nil {
			return proc, nil
		}
	}
	return nil, fmt.Errorf("all methods failed")
}
