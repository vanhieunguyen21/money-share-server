package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"money_share/pkg/auth"
	"money_share/pkg/controller"
	"money_share/pkg/database"
	"money_share/pkg/repository"
	"money_share/pkg/route"
	"net/http"
)

func main() {
	fmt.Println("Loading config...")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Cannot read config, exiting...")
		return
	}
	auth.JWTKey = []byte(viper.GetString("JWT_KEY"))

	fmt.Println("Connecting to database...")
	db := database.Connect()
	controller.UserRepository = repository.NewUserRepository(db.DB)
	controller.GroupRepository = repository.NewGroupRepository(db.DB)
	controller.MemberRepository = repository.NewMemberRepository(db.DB)
	controller.ExpenseRepository = repository.NewExpenseRepository(db.DB)

	fmt.Println("Starting server at port 8080...")
	r := mux.NewRouter()
	route.RegisterUserRoutes(r)
	route.RegisterGroupRoutes(r)
	route.RegisterMemberRoutes(r)
	route.RegisterExpenseRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
