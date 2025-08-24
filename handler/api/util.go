package api

import (
	"encoding/json"
	"net/http"
	"server/util"
)

func SendJSONResponse(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		util.Log.Errorf("failed to encode JSON response, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
