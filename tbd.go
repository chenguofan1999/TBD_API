package main

import (
	"tbd/model"
	_ "tbd/model"
	"tbd/router"
	_ "tbd/router"
)

func main() {
	testData()
	router := router.InitRouter()
	router.Run(":8011")

}

func testData() {

	model.InsertUser("Lee", "678") // 1
	model.InsertUser("Law", "678") // 2
	model.InsertUser("Jim", "678") // 3
	model.InsertUser("Tom", "678") // 4
	model.InsertUser("Bob", "678") // 5

	model.InsertFollowRelationByName("Bob", "Tom")
	model.InsertFollowRelationByName("Law", "Tom")
	model.InsertFollowRelationByName("Jim", "Tom")
	model.InsertFollowRelationByName("Lee", "Tom")
	model.InsertFollowRelationByName("Law", "Lee")
	model.InsertFollowRelationByName("Law", "Bob")
	model.InsertFollowRelationByName("Law", "Jim")
	model.InsertFollowRelationByName("Jim", "Lee")
	model.InsertFollowRelationByName("Bob", "Law")
	model.InsertFollowRelationByName("Bob", "Lee")

	model.InsertTextContent("Tom", "Hello", "Hello, world!")
	model.InsertTextContent("Law", "Hello", "Hello, world!")
	model.InsertTextContent("Jim", "Hello", "Hello, world!")
	model.InsertTextContent("Bob", "Hello", "Hello, world!")
	model.InsertTextContent("Tom", "Foo", "Bar")
	model.InsertTextContent("Law", "Foo", "Bar")

}
