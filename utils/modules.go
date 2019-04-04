package utils

import (
	"github.com/backpulse/core/constants"
)

// CheckModuleExists : loop through list of modules to check if it included
func CheckModuleExists(module string) bool {
	for _, m := range constants.Modules {
		if constants.Module(module) == m {
			return true
		}
	}
	return false
}
