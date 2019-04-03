package adminhandlers

import (
	"net/http"

	"github.com/backpulse/core/constants"
	"github.com/backpulse/core/utils"
)

// GetLanguages : return array of languages
func GetLanguages(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, "success", constants.Languages)
}
