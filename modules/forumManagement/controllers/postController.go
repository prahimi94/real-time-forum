package controller

import (
	"forum/middlewares"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	"forum/modules/forumManagement/models"
	"forum/utils"
	"net/http"
	"strconv"
	"text/template"

	userManagementControllers "forum/modules/userManagement/controllers"
	userManagementModels "forum/modules/userManagement/models"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Create a template with a function map
	tmpl, err := template.New("posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl + "posts.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func AdminReadAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	posts, err := models.ReadAllPosts()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser userManagementModels.User
		Posts     []models.Post
	}{
		LoginUser: userManagementModels.User{},
		Posts:     posts,
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("admin_posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"admin_posts.html",
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

func ReadPostsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	categoryName, errUrl := utils.ExtractFromUrl(r.URL.Path, "posts")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	filteredCategory, errCategory := models.ReadCategoryByName(categoryName)
	if errCategory != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	posts, err := models.ReadPostsByCategoryId(filteredCategory.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser            userManagementModels.User
		Posts                []models.Post
		Categories           []models.Category
		SelectedCategoryName string
	}{
		LoginUser:            userManagementModels.User{},
		Posts:                posts,
		Categories:           categories,
		SelectedCategoryName: categoryName,
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("category_posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"category_posts.html",
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

func FilterPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}
	searchTerm := r.URL.Query().Get("post_info")

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.FilterPosts(searchTerm)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
		SearchTerm string
	}{
		LoginUser:  userManagementModels.User{},
		Posts:      posts,
		Categories: categories,
		SearchTerm: searchTerm,
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("search_posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"search_posts.html",
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

func ReadMyCreatedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/myCreatedPosts/" {
		// If the URL is not exactly "/myCreatedPosts/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsByUserId(loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Posts:      posts,
		Categories: categories,
	}

	// Create a template with a function map
	tmpl, err := template.New("my_created_posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"my_created_posts.html",
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

func ReadMyLikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	if r.URL.Path != "/myLikedPosts/" {
		// If the URL is not exactly "/myLikedPosts/", respond with 404
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsLikedByUserId(loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Posts:      posts,
		Categories: categories,
	}

	// Create a template with a function map
	tmpl, err := template.New("my_liked_posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"my_liked_posts.html",
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

func ReadPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	uuid, errUrl := utils.ExtractUUIDFromUrl(r.URL.Path, "post")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(uuid, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser userManagementModels.User
		Post      models.Post
		Comments  []models.Comment
	}{
		LoginUser: loginUser,
		Post:      post,
		Comments:  nil,
	}

	if loginStatus {
		comments, err := models.ReadAllCommentsForPostByUserID(post.ID, loginUser.ID)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender.Comments = comments
	} else {
		comments, err := models.ReadAllCommentsForPost(post.ID)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender.Comments = comments
	}

	// Create a template with a function map
	tmpl, err := template.New("post_details.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"post_details.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Execute template with data
	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	// Get loginUser from context
	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"new_post.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
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

func SubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	title := utils.SanitizeInput(r.FormValue("title"))
	description := utils.SanitizeInput(r.FormValue("description"))
	categories := r.Form["categories"]
	if len(title) == 0 || len(description) == 0 || len(categories) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	// Retrieve all uploaded files
	files := r.MultipartForm.File["postFiles"]

	uploadedFiles := make(map[string]string)

	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFile, err := utils.FileUpload(file, handler)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFile
	}

	post := &models.Post{
		Title:       title,
		Description: description,
		UserId:      loginUser.ID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
			return
		}
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertPost(post, categoryIds, uploadedFiles)
	if insertError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	uuid, errUrl := utils.ExtractUUIDFromUrl(r.URL.Path, "editPost")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(uuid, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Post       models.Post
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Post:       post,
		Categories: categories,
	}

	tmpl, err := template.ParseFiles(
		publicUrl+"edit_post.html",
		publicUrl+"templates/header.html",
		publicUrl+"templates/navbar.html",
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

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("id")
	uuid := utils.SanitizeInput(r.FormValue("uuid"))
	title := utils.SanitizeInput(r.FormValue("title"))
	description := utils.SanitizeInput(r.FormValue("description"))
	categories := r.Form["categories"]

	if len(idStr) == 0 || len(title) == 0 || len(description) == 0 || len(categories) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Retrieve all uploaded files
	files := r.MultipartForm.File["postFiles"]

	uploadedFiles := make(map[string]string)

	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFile, err := utils.FileUpload(file, handler)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFile
	}

	post := &models.Post{
		ID:          id,
		Title:       title,
		Description: description,
		UserId:      loginUser.ID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	}

	// Update a record while checking duplicates
	updateError := models.UpdatePost(post, categoryIds, uploadedFiles, loginUser.ID)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+uuid, http.StatusFound)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("id")

	if len(idStr) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	post_id, err := strconv.Atoi(idStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateStatusPost(post_id, "delete", loginUser.ID)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	userManagementControllers.RedirectToIndex(w, r)
}

func AdminDeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.AdminKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("id")

	if len(idStr) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	post_id, err := strconv.Atoi(idStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateStatusPost(post_id, "delete", loginUser.ID)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	userManagementControllers.RedirectToAdminIndex(w, r)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	if !ok {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	postID := r.FormValue("post_id")
	postIDInt, _ := strconv.Atoi(postID)
	var Type string
	like := r.FormValue("like_post")
	dislike := r.FormValue("dislike_post")
	if like == "like" {
		Type = like
	} else if dislike == "dislike" {
		Type = dislike
	}

	existingLikeId, existingLikeType := models.PostHasLiked(loginUser.ID, postIDInt)

	if existingLikeId == -1 {
		post := &models.PostLike{
			Type:   Type,
			PostId: postIDInt,
			UserId: loginUser.ID,
		}
		_, insertError := models.InsertPostLike(post)
		if insertError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}
		userManagementControllers.RedirectToPrevPage(w, r)
	} else {
		updateError := models.UpdateStatusPostLike(existingLikeId, "delete", loginUser.ID)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		if existingLikeType != Type { //this is duplicated like or duplicated dislike so we should update it to disable
			post := &models.PostLike{
				Type:   Type,
				PostId: postIDInt,
				UserId: loginUser.ID,
			}
			_, insertError := models.InsertPostLike(post)
			if insertError != nil {
				errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
				return
			}
		}
		userManagementControllers.RedirectToPrevPage(w, r)
		return
	}
}
