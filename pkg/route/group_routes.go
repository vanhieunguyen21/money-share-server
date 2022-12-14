package route

import (
	"github.com/gorilla/mux"
	"money_share/pkg/controller"
	"money_share/pkg/middleware"
)

var RegisterGroupRoutes = func(router *mux.Router) {
	groupRouter := router.PathPrefix("/group").Subrouter()
	groupRouter.HandleFunc("/{groupId:[0-9]+}", controller.GetGroupById).Methods("GET")
	groupRouter.HandleFunc("/user/{userId:[0-9]+}", controller.GetGroupsByUser).Methods("GET")
	groupRouter.HandleFunc("", controller.CreateGroup).Methods("POST")
	groupRouter.HandleFunc("/{groupId:[0-9]+}", controller.UpdateGroup).Methods("PUT")
	groupRouter.HandleFunc("/{groupId:[0-9]+}", controller.DeleteGroup).Methods("DELETE")
	groupRouter.Use(middleware.Authenticate)
}