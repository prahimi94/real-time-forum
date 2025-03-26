package routes

import (
	"forum/middlewares"
	forumManagementControllers "forum/modules/forumManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	// Initialize a new router
	router := mux.NewRouter().StrictSlash(true)

	// Serve static files
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("assets/")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("assets/")))
	router.PathPrefix("/img/").Handler(http.FileServer(http.Dir("assets/")))
	router.PathPrefix("/vendor/").Handler(http.FileServer(http.Dir("assets/")))
	router.PathPrefix("/uploads/").Handler(http.FileServer(http.Dir("static/")))

	// Public routes (directly registered)
	router.HandleFunc("/", forumManagementControllers.MainPageHandler).Methods("GET")
	router.HandleFunc("/auth/", userManagementControllers.AuthHandler).Methods("GET")
	router.HandleFunc("/logout/", userManagementControllers.Logout).Methods("GET")
	router.HandleFunc("/register", userManagementControllers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", userManagementControllers.LoginHandler).Methods("POST")
	router.HandleFunc("/post/{id}", forumManagementControllers.ReadPost).Methods("GET")
	router.HandleFunc("/posts/{categoryName}", forumManagementControllers.ReadPostsByCategory).Methods("GET")
	router.HandleFunc("/filterPosts", forumManagementControllers.FilterPosts).Methods("GET")

	// Protected routes (using middleware)
	protectedRoutes := router.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middlewares.AuthMiddleware) // Apply AuthMiddleware to protectedRoutes routes
	protectedRoutes.HandleFunc("/profile", userManagementControllers.EditUser).Methods("GET")
	protectedRoutes.HandleFunc("/updateUser", userManagementControllers.UpdateUser).Methods("POST")
	protectedRoutes.HandleFunc("/newPost/", forumManagementControllers.CreatePost).Methods("GET")
	protectedRoutes.HandleFunc("/submitPost", forumManagementControllers.SubmitPost).Methods("POST")
	protectedRoutes.HandleFunc("/editPost/{id}", forumManagementControllers.EditPost).Methods("GET")
	protectedRoutes.HandleFunc("/updatePost", forumManagementControllers.UpdatePost).Methods("POST")
	protectedRoutes.HandleFunc("/deletePost", forumManagementControllers.DeletePost).Methods("POST")
	protectedRoutes.HandleFunc("/myCreatedPosts/", forumManagementControllers.ReadMyCreatedPosts).Methods("GET")
	protectedRoutes.HandleFunc("/myLikedPosts/", forumManagementControllers.ReadMyLikedPosts).Methods("GET")
	protectedRoutes.HandleFunc("/likePost", forumManagementControllers.LikePost).Methods("POST")
	protectedRoutes.HandleFunc("/likeComment", forumManagementControllers.LikeComment).Methods("POST")
	protectedRoutes.HandleFunc("/submitComment", forumManagementControllers.SubmitComment).Methods("POST")
	protectedRoutes.HandleFunc("/updateComment", forumManagementControllers.UpdateComment).Methods("POST")
	protectedRoutes.HandleFunc("/deleteComment", forumManagementControllers.DeleteComment).Methods("POST")

	// Protected routes (using middleware)
	adminRoutes := router.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middlewares.AdminMiddleware) // Apply AdminMiddleware to admin routes
	adminRoutes.HandleFunc("/", forumManagementControllers.AdminMainPageHandler).Methods("GET")
	// adminRoutes.HandleFunc("/dashboard", forumManagementControllers.AdminDashboardHandler).Methods("GET")
	adminRoutes.HandleFunc("/users", userManagementControllers.AdminReadAllUsers).Methods("GET")
	adminRoutes.HandleFunc("/posts", forumManagementControllers.AdminReadAllPosts).Methods("GET")
	adminRoutes.HandleFunc("/deletePost", forumManagementControllers.AdminDeletePost).Methods("POST")
	// adminRoutes.HandleFunc("/posts", adminControllers.ManagePosts).Methods("GET")
	// adminRoutes.HandleFunc("/deleteUser", adminControllers.DeleteUser).Methods("POST")

	return router
}
