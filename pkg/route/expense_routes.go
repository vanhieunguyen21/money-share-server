package route

import (
	"github.com/gorilla/mux"
	"money_share/pkg/controller"
	"money_share/pkg/middleware"
)

var RegisterExpenseRoutes = func(router *mux.Router) {
	expenseRouter := router.PathPrefix("/expense").Subrouter()
	expenseRouter.HandleFunc("/{expenseId:[0-9]+}", controller.GetExpenseByID).Methods("GET")
	expenseRouter.HandleFunc("/group/{groupId:[0-9]+}", controller.GetExpensesByGroup).Methods("GET")
	expenseRouter.HandleFunc("", controller.GetExpensesByMember).Methods("GET")
	expenseRouter.HandleFunc("", controller.CreateExpense).Methods("POST")
	expenseRouter.HandleFunc("/{expenseId:[0-9]+}", controller.UpdateExpense).Methods("PUT")
	expenseRouter.HandleFunc("/{expenseId:[0-9]+}", controller.DeleteExpense).Methods("DELETE")
	expenseRouter.Use(middleware.Authenticate)

}
