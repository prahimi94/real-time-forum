package controller

import (
	"fmt"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	"forum/modules/forumManagement/models"
	"forum/utils"
	"net/http"
	"text/template"

	userManagementControllers "forum/modules/userManagement/controllers"
	userManagementModels "forum/modules/userManagement/models"

	_ "github.com/mattn/go-sqlite3"
)

const publicUrl = "modules/forumManagement/views/"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse template
	tmpl, err := template.ParseFiles(publicUrl + "index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/" {
		// If the URL is not exactly "/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		categories, err := models.ReadAllCategories()
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		posts, err := models.ReadAllPosts(loginUser.ID)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender := struct {
			LoginUser  userManagementModels.User
			Posts      []models.Post
			Categories []models.Category
		}{
			LoginUser:  userManagementModels.User{},
			Posts:      posts,
			Categories: categories,
		}

		if loginStatus {
			if loginUser.Type == "admin" {
				userManagementControllers.RedirectToAdminIndex(w, r)
				return
			}
			data_obj_sender.LoginUser = loginUser
		}

		// jsonData, err := json.Marshal(data_obj_sender)
		// if err != nil {
		// 	http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		// 	return
		// }

		// w.Header().Set("Content-Type", "application/json")
		// w.Write(jsonData) // Manually writing JSON to response

		// w.Header().Set("Content-Type", "application/json")
		// if err := json.NewEncoder(w).Encode(data_obj_sender); err != nil {
		// 	http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		// }

		// Create a template with a function map
		tmpl, err := template.New("index.html").Funcs(template.FuncMap{
			"formatDate": utils.FormatDate, // Register function globally
		}).ParseFiles(
			publicUrl + "index.html",
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
	} else {
		tmpl, err := template.New("index.html").ParseFiles(
			publicUrl + "index.html",
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

}

func AdminMainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/admin/" {
		// If the URL is not exactly "/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	categories, err := models.AdminReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadAllPosts(loginUser.ID)
	if err != nil {
		fmt.Println(1)
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	comments, err := models.ReadAllComments()
	if err != nil {
		fmt.Println(2)
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	users, err := userManagementModels.ReadAllUsers()
	if err != nil {
		fmt.Println(3)
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	postLikes, err := models.ReadAllPostsLikes()
	if err != nil {
		fmt.Println(3)
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Comments   []models.Comment
		Users      []userManagementModels.User
		PostLikes  []models.PostLike
		Categories []models.Category
	}{
		LoginUser:  userManagementModels.User{},
		Posts:      posts,
		Comments:   comments,
		Users:      users,
		PostLikes:  postLikes,
		Categories: categories,
	}

	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("admin_dashboard.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"admin_dashboard.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
		publicUrl+"templates/hero.html",
		publicUrl+"templates/posts.html",
		publicUrl+"templates/footer.html",
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
