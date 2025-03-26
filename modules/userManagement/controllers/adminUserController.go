package controller

import (
	"fmt"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	"forum/modules/userManagement/models"
	"forum/utils"
	"net/http"
	"text/template"
)

// AdminReadAllUsers retrieves and displays all users for admin purposes.
func AdminReadAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/admin/users" {
		// If the URL is not exactly "/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	users, err := models.ReadAllUsers()
	if err != nil {
		fmt.Println(3)
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser models.User
		Users     []models.User
	}{
		LoginUser: models.User{},
		Users:     users,
	}

	loginStatus, loginUser, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("admin_users.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"admin_users.html",
		forumPublicUrl+"templates/footer.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

// AdminUpdateUser handles updating a user's status or details.
func AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract parameters, validate admin rights, and update using the User model.
	fmt.Fprintln(w, "Admin: User updated")
}

// AdminDeleteUser handles user deletion.
func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Use the model to delete the user.
	fmt.Fprintln(w, "Admin: User deleted")
}

func RedirectToAdminIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/", http.StatusFound)
}
