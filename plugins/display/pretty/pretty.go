// plugin.go
package main

import (
	"github.com/k0kubun/pp"
)

// OutputDisplayModuleImpl is an instance of the OutputDisplayModule interface.
var DisplayPretty OutputDisplayModule

type OutputDisplayModule struct{}

// Display is a method of OutputDisplayModule that calls the Display method of DisplayPlugin.
func (m OutputDisplayModule) Display(output interface{}) error {
	// Call the Display method of the plugin
	pp.Println(output)
	return nil
}
