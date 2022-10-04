package route

import (
	"github.com/gorilla/mux"
	"money_share/pkg/controller"
)

var RegisterExpenseRoutes = func(router *mux.Router) {
	router.HandleFunc("/expense/{expenseId:[0-9]+}", controller.GetExpenseByID).Methods("GET")
	router.HandleFunc("/expense/group/{groupId:[0-9]+}", controller.GetExpensesByGroup).Methods("GET")
	router.HandleFunc("/expense", controller.GetExpensesByMember).Methods("GET")
	router.HandleFunc("/expense", controller.CreateExpense).Methods("POST")
	router.HandleFunc("/expense/{expenseId:[0-9]+}", controller.UpdateExpense).Methods("PUT")
	router.HandleFunc("/expense/{expenseId:[0-9]+}", controller.DeleteExpense).Methods("DELETE")
}
