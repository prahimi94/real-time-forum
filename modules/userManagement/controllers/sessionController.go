package controller

import (
	"fmt"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	"forum/utils"

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

	loginStatus, loginUser, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
	} else {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	data_obj_sender := struct {
		LoginUser userManagementModels.User
		Active    bool
	}{
		LoginUser: loginUser,
		Active:    isActive,
	}

	// Respond with the session status
	res := utils.Result{
		Success: true,
		Message: "",
		Data:    data_obj_sender,
	}
	utils.ReturnJson(w, res)
}
