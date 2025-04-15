package controller

import (
	"encoding/json"
	"fmt"
	errorManagementControllers "forum/modules/errorManagement/controllers"

	userManagementModels "forum/modules/userManagement/models"
	"forum/utils"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const publicUrl = "modules/userManagement/views/"
const forumPublicUrl = "modules/forumManagement/views/"

var OnlineUsers = make(map[*websocket.Conn]string) // Map of online users (connected to WS) to usernames

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
	// if loginStatus {
	// 	RedirectToIndex(w, r)
	// 	return
	// }
	if loginStatus {
		res := utils.Result{
			Success: true,
			Message: "You are logged in",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	username := utils.SanitizeInput(r.FormValue("username"))
	firstname := utils.SanitizeInput(r.FormValue("firstname"))
	lastname := utils.SanitizeInput(r.FormValue("lastname"))
	gender := utils.SanitizeInput(r.FormValue("gender"))
	age := utils.SanitizeInput(r.FormValue("age"))
	email := utils.SanitizeInput(r.FormValue("email"))
	password := utils.SanitizeInput(r.FormValue("password"))
	if len(username) == 0 || len(firstname) == 0 || len(lastname) == 0 || len(gender) == 0 || len(age) == 0 || len(email) == 0 || len(password) == 0 {
		// errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		// renderAuthPage(w, "Username, email and password are required.")
		res := utils.Result{
			Success:    false,
			Message:    "Username, firstname, lastname, gender, age, email and password are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		// renderAuthPage(w, "Invalid email address!")
		res := utils.Result{
			Success:    false,
			Message:    "Invalid email address!",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		// renderAuthPage(w, "Invalid age format!")
		res := utils.Result{
			Success:    false,
			Message:    "Invalid age format!",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	newUser := &userManagementModels.User{
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Gender:    gender,
		Age:       ageInt,
		Email:     email,
		Password:  string(hashedPassword),
	}

	// Insert a record while checking duplicates
	userId, insertError := userManagementModels.InsertUser(newUser)
	if insertError != nil {
		if insertError.Error() == "duplicateEmail" {
			// renderAuthPage(w, "User with this email already exists!")
			res := utils.Result{
				Success:    false,
				Message:    "User with this email already exists!",
				HttpStatus: http.StatusOK,
				Data:       nil,
			}
			utils.ReturnJson(w, res)
			return
		} else if insertError.Error() == "duplicateUsername" {
			// renderAuthPage(w, "User with this username already exists!")
			res := utils.Result{
				Success:    false,
				Message:    "User with this username already exists!",
				HttpStatus: http.StatusOK,
				Data:       nil,
			}
			utils.ReturnJson(w, res)
			return
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		}
		return
	}

	sessionGenerator(w, r, userId)

	res := utils.Result{
		Success: true,
		Message: "Logged in successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
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
	// if loginStatus {
	// 	RedirectToIndex(w, r)
	// 	return
	// }
	if loginStatus {
		res := utils.Result{
			Success: true,
			Message: "You are logged in",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	username := utils.SanitizeInput(r.FormValue("username"))
	password := utils.SanitizeInput(r.FormValue("password"))
	if len(username) == 0 || len(password) == 0 {
		// errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		// renderAuthPage(w, "Username and password are required.")
		res := utils.Result{
			Success:    false,
			Message:    "Username and password are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	// Insert a record while checking duplicates
	authStatus, userId, authError := userManagementModels.AuthenticateUser(username, password)
	if authError != nil {
		// errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		// renderAuthPage(w, authError.Error())
		res := utils.Result{
			Success:    false,
			Message:    authError.Error(),
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	if authStatus {
		sessionGenerator(w, r, userId)
	}

	res := utils.Result{
		Success: true,
		Message: "Logged in successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
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
	UserSetCookie(w, session.SessionToken, session.ExpiresAt)
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
	loginStatus, loggedInUser, sessionToken, checkLoginError := CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	if !loginStatus {
		// RedirectToIndex(w, r)
		res := utils.Result{
			Success: true,
			Message: "You are not logged in",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	err := userManagementModels.DeleteSession(sessionToken)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	deleteCookie(w, "session_token") // Deleting a cookie named "session_token"
	SocketLogoutHandler(w, r, loggedInUser.Username)
	// delete ws conn
	//for client := range forumManagementControllers.OnlineUsers {
	//if loggedInUser.Username == client {
	//}
	//fmt.Println(client)
	/* forumManagementControllers.Mutex.Lock()
	delete(forumManagementControllers.OnlineUsers, forumManagementControllers.conn)
	forumManagementControllers.UpdateOnlineUsers() */
	//}
	//forumManagementControllers.Mutex.Unlock()

	// RedirectToIndex(w, r)
	res := utils.Result{
		Success: true,
		Message: "Logged out successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
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
	if !loginStatus {
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
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	const maxUploadSize = 2 << 20 // 2 MB

	// Limit the request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	firstname := utils.SanitizeInput(r.FormValue("firstname"))
	lastname := utils.SanitizeInput(r.FormValue("lastname"))
	age := utils.SanitizeInput(r.FormValue("age"))
	gender := utils.SanitizeInput(r.FormValue("gender"))

	if len(firstname) == 0 || len(lastname) == 0 || len(age) == 0 || len(gender) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "firstname, lastname, gender and age are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		// renderAuthPage(w, "Invalid age format!")
		res := utils.Result{
			Success:    false,
			Message:    "Invalid age format!",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	profile_photo_file, handler, err := r.FormFile("profile_photo")
	if err != nil {
		// "File is missing"

		user := &userManagementModels.User{
			ID:           loginUser.ID,
			Firstname:    firstname,
			Lastname:     lastname,
			Gender:       gender,
			Age:          ageInt,
			ProfilePhoto: "",
		}

		// Update a record while checking duplicates
		updateError := userManagementModels.UpdateUser(user)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		res := utils.Result{
			Success: true,
			Message: "Profile updated successfully",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	} else {
		defer profile_photo_file.Close()

		profile_photo := ""
		if handler.Size != 0 {
			// Extra safety: check file size from the header
			if handler.Size > maxUploadSize {
				// "File is too large or missing"
				errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
				return
			}

			// Call your file upload function
			profile_photo, err = utils.FileUpload(profile_photo_file, handler)
			if err != nil {
				errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
				return
			}
		}

		user := &userManagementModels.User{
			ID:           loginUser.ID,
			Firstname:    firstname,
			Lastname:     lastname,
			Gender:       gender,
			Age:          ageInt,
			ProfilePhoto: profile_photo,
		}

		// Update a record while checking duplicates
		updateError := userManagementModels.UpdateUser(user)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		res := utils.Result{
			Success: true,
			Message: "Profile updated successfully",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}

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
		Name:     cookieName,
		Value:    "",              // Optional but recommended
		Expires:  time.Unix(0, 0), // Set expiration to a past date
		MaxAge:   -1,              // Ensure immediate removal
		Path:     "/",             // Must match the original cookie path
		HttpOnly: true,
		Secure:   false,
	})
}

func UserSetCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
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

// Helper function to broadcast the list of online users
func UpdateOnlineUsers() {
	usernames := make([]string, 0, len(OnlineUsers))
	for _, username := range OnlineUsers {
		usernames = append(usernames, username)
	}

	// Encode the list of usernames as JSON
	userListJSON, err := json.Marshal(usernames)
	if err != nil {
		fmt.Println("Error encoding online users:", err)
		return
	}

	// Send the list to all online clients
	for client := range OnlineUsers {
		err := client.WriteMessage(websocket.TextMessage, userListJSON)
		if err != nil {
			client.Close()
			delete(OnlineUsers, client)
		}
	}
}

func SocketLogoutHandler(w http.ResponseWriter, r *http.Request, userName string) {
	for clientConn, clientUserName := range OnlineUsers {
		if clientUserName == userName {
			delete(OnlineUsers, clientConn)
			UpdateOnlineUsers()
		}
		clientConn.WriteJSON(map[string]string{
			"type":    "logout",
			"message": "You have been logged out.",
		})
		clientConn.Close() // close their websocket
	}
}
