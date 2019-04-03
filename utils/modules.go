package utils

import (
	"github.com/backpulse/core/constants"
)

func CheckModuleExists(module string) bool {
	for _, m := range constants.Modules {
		if constants.Module(module) == m {
			return true
		}
	}
	return false
}
