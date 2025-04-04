package routes

import (
	forumManagementControllers "forum/modules/forumManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
	"net/http"
)

func SetupRoutes() {
	http.Handle("/css/", http.FileServer(http.Dir("assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	http.Handle("/img/", http.FileServer(http.Dir("assets/")))
	http.Handle("/uploads/", http.FileServer(http.Dir("static/")))

	// Register route handlers
	http.HandleFunc("/", forumManagementControllers.MainPageHandler)
	// http.HandleFunc("/", forumManagementControllers.IndexHandler)

	//FOR CHAT
	http.HandleFunc("/ws", forumManagementControllers.WsHandler)
	http.HandleFunc("/api/connected-users", forumManagementControllers.OnlineUsersHandler) //GET USERS CONNECTED TO WS

	//FOR SESSIONS
	http.HandleFunc("/api/check-session", userManagementControllers.CheckSessionHandler)

	//TO GET LOGGED-IN USERS = active session
	http.HandleFunc("/api/loggedin-users", userManagementControllers.LoggedInUsersHandler)

	// http.HandleFunc("/home/", forumManagementControllers.HomePageHandler)
	http.HandleFunc("/auth/", userManagementControllers.AuthHandler)
	http.HandleFunc("/logout/", userManagementControllers.Logout)
	http.HandleFunc("/register", userManagementControllers.RegisterHandler) /*post method*/
	http.HandleFunc("/login", userManagementControllers.LoginHandler)       /*post method*/

	http.HandleFunc("/api/categories/", forumManagementControllers.ReadAllCategories)
	http.HandleFunc("/api/allPosts/", forumManagementControllers.ReadAllPosts)
	// http.HandleFunc("/newPost/", forumManagementControllers.CreatePost)
	http.HandleFunc("/api/submitPost", forumManagementControllers.SubmitPost) /*post method*/
	http.HandleFunc("/api/myCreatedPosts/", forumManagementControllers.ReadMyCreatedPosts)
	http.HandleFunc("/api/myLikedPosts/", forumManagementControllers.ReadMyLikedPosts)

	// router.HandleFunc("/post/{id}", forumManagementControllers.ReadPost).Methods("GET")
	http.HandleFunc("/api/post/", forumManagementControllers.ReadPost)
	// router.HandleFunc("/posts/{categoryName}", forumManagementControllers.ReadPostsByCategory).Methods("GET")
	http.HandleFunc("/api/posts/", forumManagementControllers.ReadPostsByCategory)
	http.HandleFunc("/api/filterPosts/", forumManagementControllers.FilterPosts)
	http.HandleFunc("/api/likePost", forumManagementControllers.LikePost)
	// protectedRoutes.HandleFunc("/editPost/{id}", forumManagementControllers.EditPost).Methods("GET")
	http.HandleFunc("/editPost/", forumManagementControllers.EditPost)
	http.HandleFunc("/api/updatePost", forumManagementControllers.UpdatePost) /*post method*/
	http.HandleFunc("/api/deletePost", forumManagementControllers.DeletePost) /*post method*/

	http.HandleFunc("/api/likeComment", forumManagementControllers.LikeComment)
	http.HandleFunc("/api/submitComment", forumManagementControllers.SubmitComment) /*post method*/
	http.HandleFunc("/api/updateComment", forumManagementControllers.UpdateComment) /*post method*/
	http.HandleFunc("/api/deleteComment", forumManagementControllers.DeleteComment) /*post method*/

	// // Protected routes (using middleware)
	// protectedRoutes := router.PathPrefix("/").Subrouter()
	// protectedRoutes.Use(middlewares.AuthMiddleware) // Apply AuthMiddleware to protectedRoutes routes
	// protectedRoutes.HandleFunc("/profile", userManagementControllers.EditUser).Methods("GET")
	// protectedRoutes.HandleFunc("/updateUser", userManagementControllers.UpdateUser).Methods("POST")

	// // Protected routes (using middleware)
	// adminRoutes := router.PathPrefix("/admin").Subrouter()
	// adminRoutes.Use(middlewares.AdminMiddleware) // Apply AdminMiddleware to admin routes
	// adminRoutes.HandleFunc("/", forumManagementControllers.AdminMainPageHandler).Methods("GET")
	// // adminRoutes.HandleFunc("/dashboard", forumManagementControllers.AdminDashboardHandler).Methods("GET")
	// adminRoutes.HandleFunc("/users", userManagementControllers.AdminReadAllUsers).Methods("GET")
	// adminRoutes.HandleFunc("/posts", forumManagementControllers.AdminReadAllPosts).Methods("GET")
	// adminRoutes.HandleFunc("/deletePost", forumManagementControllers.AdminDeletePost).Methods("POST")
	// // adminRoutes.HandleFunc("/posts", adminControllers.ManagePosts).Methods("GET")
	// // adminRoutes.HandleFunc("/deleteUser", adminControllers.DeleteUser).Methods("POST")
}
