package version

import "fmt"

var Major int = 1
var Minor int = 0
var Patch int = 0

func Formatted() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}
