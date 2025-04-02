package controller

import (
	"encoding/json"
	"fmt"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	userManagementModels "forum/modules/userManagement/models"
	"forum/utils"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const publicUrl = "modules/userManagement/views/"
const forumPublicUrl = "modules/forumManagement/views/"

//var u1 = uuid.Must(uuid.NewV4())

type AuthPageErrorData struct {
	ErrorMessage string
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, _, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		RedirectToIndex(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		publicUrl + "authPage.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, _, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		RedirectToIndex(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	username := utils.SanitizeInput(r.FormValue("username"))
	email := utils.SanitizeInput(r.FormValue("email"))
	password := utils.SanitizeInput(r.FormValue("password"))
	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		// errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		renderAuthPage(w, "Username, email and password are required.")
		return
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		renderAuthPage(w, "Invalid email address!")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	newUser := &userManagementModels.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Insert a record while checking duplicates
	userId, insertError := userManagementModels.InsertUser(newUser)
	if insertError != nil {
		if insertError.Error() == "duplicateEmail" {
			renderAuthPage(w, "User with this email already exists!")
			return
		} else if insertError.Error() == "duplicateUsername" {
			renderAuthPage(w, "User with this username already exists!")
			return
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		}
		return
	}

	sessionGenerator(w, r, userId)

	RedirectToIndex(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, _, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		RedirectToIndex(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	username := utils.SanitizeInput(r.FormValue("username"))
	password := utils.SanitizeInput(r.FormValue("password"))
	if len(username) == 0 || len(password) == 0 {
		// errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		// return
		renderAuthPage(w, "Username and password are required.")
		return
	}

	// Insert a record while checking duplicates
	authStatus, userId, authError := userManagementModels.AuthenticateUser(username, password)
	if authError != nil {
		// errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		renderAuthPage(w, authError.Error())
		return
	}
	if authStatus {
		sessionGenerator(w, r, userId)
	}

	RedirectToIndex(w, r)
}

// Render the login page with an optional error message
func renderAuthPage(w http.ResponseWriter, errorMsg string) {
	tmpl := template.Must(template.ParseFiles(publicUrl + "authPage.html"))
	tmpl.Execute(w, AuthPageErrorData{ErrorMessage: errorMsg})
}

func sessionGenerator(w http.ResponseWriter, r *http.Request, userId int) {
	session := &userManagementModels.Session{
		UserId: userId,
	}
	session, insertError := userManagementModels.InsertSession(session)
	if insertError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	SetCookie(w, session.SessionToken, session.ExpiresAt)
	// Set the session token in a cookie

}

// Middleware to check for valid user session in cookie
func CheckLogin(w http.ResponseWriter, r *http.Request) (bool, userManagementModels.User, string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, userManagementModels.User{}, "", nil
	}

	sessionToken := cookie.Value
	user, expirationTime, selectError := userManagementModels.SelectSession(sessionToken)
	if selectError != nil {
		if selectError.Error() == "sql: no rows in result set" {
			deleteCookie(w, "session_token")
			return false, userManagementModels.User{}, "", nil
		} else {
			return false, userManagementModels.User{}, "", selectError
		}
	}

	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		// Cookie expired, redirect to login
		return false, userManagementModels.User{}, "", nil
	}

	return true, user, sessionToken, nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
	loginStatus, _, sessionToken, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	if !loginStatus {
		RedirectToIndex(w, r)
		return
	}

	err := userManagementModels.DeleteSession(sessionToken)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	deleteCookie(w, "session_token") // Deleting a cookie named "session_token"
	RedirectToIndex(w, r)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
		// return
	} else {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	data_obj_sender := struct {
		LoginUser userManagementModels.User
	}{
		LoginUser: loginUser,
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"edit_user.html",
		forumPublicUrl+"templates/header.html",
		forumPublicUrl+"templates/navbar.html",
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		fmt.Println("logged in userid is: ", loginUser.ID)
		// return
	} else {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	const maxUploadSize = 2 << 20 // 2 MB

	// Limit the request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	name := utils.SanitizeInput(r.FormValue("name"))

	if len(name) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	profile_photo_file, handler, err := r.FormFile("profile_photo")
	if err != nil {
		// "File is too large or missing"
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	defer profile_photo_file.Close()

	// Extra safety: check file size from the header
	if handler.Size > maxUploadSize {
		// "File is too large or missing"
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	// Call your file upload function
	profile_photo, err := utils.FileUpload(profile_photo_file, handler)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	user := &userManagementModels.User{
		ID:           loginUser.ID,
		Name:         name,
		ProfilePhoto: profile_photo,
	}

	// Update a record while checking duplicates
	updateError := userManagementModels.UpdateUser(user)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	RedirectToIndex(w, r)
}

func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func RedirectToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home/", http.StatusFound)
}

func RedirectToPrevPage(w http.ResponseWriter, r *http.Request) {
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		referrer = "/"
	}

	// Redirect back to the original page to reload it
	http.Redirect(w, r, referrer, http.StatusSeeOther)
}

func deleteCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   "",              // Optional but recommended
		Expires: time.Unix(0, 0), // Set expiration to a past date
		MaxAge:  -1,              // Ensure immediate removal
		Path:    "/",             // Must match the original cookie path
	})
}

func SetCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
	})
}

func LoggedInUsersHandler(w http.ResponseWriter, r *http.Request) {
	usernames, err := userManagementModels.GetActiveSessionUsernames(r)
	if err != nil {
		http.Error(w, "Failed to fetch logged-in users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usernames)
}
