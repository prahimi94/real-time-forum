package middlewares

import (
	"context"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	userManagementControllers "forum/modules/userManagement/controllers"
	"net/http"
)

// AdminContextKey is used to store user info in the request context
type AdminContextKey string

const AdminKey AdminContextKey = "admin"

// AuthMiddleware ensures that only authenticated users can access certain routes
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
		if checkLoginError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		if loginStatus && loginUser.Type == "admin" {
			// Store admin user in context
			ctx := context.WithValue(r.Context(), AdminKey, loginUser)

			// Pass request with updated context to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
			// return
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.ForbiddenError)
			return
		}
	})
}
