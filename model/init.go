package model

func init() {
	Connect()
	// clearDB() // remove this when release
	CreateUserTableIfNotExists()
	CreateFollowTableIfNotExists()
	CreateContentTableIfNotExists()
	CreateImageTableIfNotExists()
	CreateCommentTableIfNotExists()
	CreateLikeTableIfNotExists()
}

func clearDB() {
	DB.Exec("drop table likes")
	DB.Exec("drop table comments")
	DB.Exec("drop table images")
	DB.Exec("drop table contents")
	DB.Exec("drop table follow")
	DB.Exec("drop table users")
}

func testData() {

	InsertUser("Lee", "678") // 1
	InsertUser("Law", "678") // 2
	InsertUser("Jim", "678") // 3
	InsertUser("Tom", "678") // 4
	InsertUser("Bob", "678") // 5

	InsertFollowRelationByName("Bob", "Tom")
	InsertFollowRelationByName("Law", "Tom")
	InsertFollowRelationByName("Jim", "Tom")
	InsertFollowRelationByName("Lee", "Tom")
	InsertFollowRelationByName("Law", "Lee")
	InsertFollowRelationByName("Law", "Bob")
	InsertFollowRelationByName("Law", "Jim")
	InsertFollowRelationByName("Jim", "Lee")
	InsertFollowRelationByName("Bob", "Law")
	InsertFollowRelationByName("Bob", "Lee")

	InsertContent("Tom", "Hello", "Hello, world!", []string{"static/images/Yosemite-Color-Block.png"})
	InsertContent("Law", "Hello", "Hello, world!", []string{"static/images/Yosemite.png", "static/images/Backgrounds.png"})
	InsertContent("Jim", "Hello", "Hello, world!", []string{})
	InsertContent("Bob", "Hello", "Hello, world!", []string{})
	InsertContent("Tom", "Foo", "Bar", []string{})
	InsertContent("Law", "Foo", "Bar", []string{})

	InsertComment("Tom", 2, "nice")  // commentID = 1
	InsertComment("Jim", 1, "great") // commentID = 2
	InsertComment("Bob", 3, "fair")  // commentID = 3
	InsertComment("Law", 2, "yep")   // commentID = 4

	InsertReply("Jim", 1, "thx") // commentID = 5
	InsertReply("Bob", 1, "na")  // commentID = 6
	InsertReply("Law", 1, "na")  // commentID = 7
	InsertReply("Tom", 1, "na")  // commentID = 8

	InsertReply("Tom", 5, "na") // commentID = 8
	InsertReply("Jim", 5, "na") // commentID = 8
}
