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
	router.Run(":8012")
	// fmt.Println(model.QueryReplyWithID(5))
	// fmt.Println(model.QueryAllCommentsWithContentID(2))
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

	model.InsertContent("Tom", "Hello", "Hello, world!", []string{"static/images/Yosemite-Color-Block.png"})
	model.InsertContent("Law", "Hello", "Hello, world!", []string{"static/images/Yosemite.png", "static/images/Backgrounds.png"})
	model.InsertContent("Jim", "Hello", "Hello, world!", []string{})
	model.InsertContent("Bob", "Hello", "Hello, world!", []string{})
	model.InsertContent("Tom", "Foo", "Bar", []string{})
	model.InsertContent("Law", "Foo", "Bar", []string{})

	model.InsertComment("Tom", 2, "nice")  // commentID = 1
	model.InsertComment("Jim", 1, "great") // commentID = 2
	model.InsertComment("Bob", 3, "fair")  // commentID = 3
	model.InsertComment("Law", 2, "yep")   // commentID = 4

	model.InsertReply("Jim", 1, "thx") // commentID = 5
	model.InsertReply("Bob", 1, "na")  // commentID = 6
	model.InsertReply("Law", 1, "na")  // commentID = 7
	model.InsertReply("Tom", 1, "na")  // commentID = 8

	model.InsertReply("Tom", 5, "na") // commentID = 8
	model.InsertReply("Jim", 5, "na") // commentID = 8
}
