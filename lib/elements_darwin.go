// build +darwin

package lib

import (
	_ "github.com/jtopjian/elements/utils"
)

// GetSystemElements is the main function to call to collect all elements on the system.
// You shouldn't need to interact with anything other than this function.
func (e *Elements) GetSystemElements() *SystemElements {
	elements := new(SystemElements)
	elements.Architecture = "darwin"
	return elements
}

type SystemElements struct {
	Architecture string
}
