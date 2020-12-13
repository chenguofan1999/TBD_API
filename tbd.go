package main

import (
	"tbd/model"
	_ "tbd/model"
	"tbd/router"
	_ "tbd/router"
)

func main() {

	router := router.InitRouter()
	router.Run(":8011")

}

func testData() {

	model.InsertUser(model.User{
		Username: "Lee",
		Password: "678",
	})

	model.InsertUser(model.User{
		Username: "Law",
		Password: "678",
	})

	model.InsertUser(model.User{
		Username: "Jim",
		Password: "678",
	})

	model.InsertUser(model.User{
		Username: "Tom",
		Password: "678",
	})

	model.InsertUser(model.User{
		Username: "Bob",
		Password: "678",
	})

	model.InsertFollowRelation(1, 2)
	model.InsertFollowRelation(1, 3)
	model.InsertFollowRelation(2, 3)
	model.InsertFollowRelation(4, 3)
	model.InsertFollowRelation(2, 1)
	model.InsertFollowRelation(1, 4)

	model.InsertTextContent("Tom", "Hello", "Hello, world!")
	model.InsertTextContent("Law", "Hello", "Hello, world!")
	model.InsertTextContent("Jim", "Hello", "Hello, world!")
	model.InsertTextContent("Bob", "Hello", "Hello, world!")
	model.InsertTextContent("Tom", "Foo", "Bar")
	model.InsertTextContent("Law", "Foo", "Bar")

}
