package controller

import (
	"encoding/json"
	userManagementModels "forum/modules/userManagement/models"
	"net/http"
)

// CheckSessionHandler checks if the user's session is active
func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session token from cookies
	cookie, err := r.Cookie("session_token")
	if err != nil {
		//http.Error(w, "Session token not found", http.StatusUnauthorized)
		return
	}

	// Check if the session is active
	isActive, err := userManagementModels.IsSessionActive(cookie.Value)
	if err != nil {
		//http.Error(w, "Error checking session status", http.StatusInternalServerError)
		return
	}

	// Respond with the session status
	response := map[string]bool{"active": isActive}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
