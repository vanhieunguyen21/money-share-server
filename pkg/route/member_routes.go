package route

import (
	"github.com/gorilla/mux"
	"money_share/pkg/controller"
)

var RegisterMemberRoutes = func(router *mux.Router) {
	router.HandleFunc("/member", controller.GetMemberByID).Methods("GET")
	router.HandleFunc("/member/group/{groupId:[0-9]+}", controller.GetMembersOfGroup).Methods("GET")
	router.HandleFunc("/member", controller.AddMemberToGroup).Methods("POST")
	router.HandleFunc("/member", controller.RemoveMemberFromGroup).Methods("DELETE")
}
