package route

import (
	"github.com/gorilla/mux"
	"money_share/pkg/controller"
	"money_share/pkg/middleware"
)

var RegisterUserRoutes = func(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", controller.Login).Methods("POST")
	userRouter.HandleFunc("/{userId:[0-9]+}", controller.GetUserByID).Methods("GET")
	userRouter.HandleFunc("/checkUsername/{username}", controller.CheckUsername).Methods("GET")
	userRouter.HandleFunc("/register", controller.Register).Methods("POST")
	userRouter.HandleFunc("/profileImage/{fileName}", controller.GetUserProfileImage).Methods("GET")

	// Authentication required routes
	authSub := userRouter.PathPrefix("/auth").Subrouter()
	authSub.HandleFunc("/{userId:[0-9]+}", controller.UpdateUser).Methods("PUT")
	authSub.HandleFunc("/{userId:[0-9]+}/profileImage", controller.UploadUserProfileImage).Methods("POST")
	authSub.HandleFunc("/{userId:[0-9]+}", controller.DeleteUser).Methods("DELETE")
	authSub.Use(middleware.Authenticate)
}
