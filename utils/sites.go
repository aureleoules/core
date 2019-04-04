package utils

import (
	"github.com/backpulse/core/models"
)

// HasPremiumFeatures : i don't know
func HasPremiumFeatures(site models.Site, size float64) bool {
	//TODO: Look into this function
	if len(site.Collaborators) > 1 {
		return true
	}
	if site.TotalSize > 500 {
		return true
	}
	return false
}
