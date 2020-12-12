package main

import (
	_ "tbd/model"
	"tbd/router"
)

func main() {
	// newUser := model.User{
	// 	Username: "Tom",
	// 	Password: "678",
	// }

	// model.InsertUser(newUser)
	// model.QueryWithName("Joe")

	router := router.InitRouter()
	router.Run(":8011")
}
