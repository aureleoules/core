package utils

import (
	"github.com/backpulse/core/models"
)

func HasPremiumFeatures(site models.Site, size float64) bool {
	if len(site.Collaborators) > 1 {
		return true
	}
	if site.TotalSize > 500 {
		return true
	}
	return false
}
