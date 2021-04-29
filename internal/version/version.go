package version

import "fmt"

// When updating this version number also update installer.nsi
var Major int = 1
var Minor int = 0
var Patch int = 0

var Protocol int = 1

func Formatted() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}
