package model

func init() {
	Connect()
	clearDB() // remove this when release
	CreateUserTableIfNotExists()
	CreateFollowTableIfNotExists()
	CreateContentTableIfNotExists()
	CreateImageTableIfNotExists()
	CreateCommentTableIfNotExists()
}

func clearDB() {
	DB.Exec("drop table comments")
	DB.Exec("drop table images")
	DB.Exec("drop table contents")
	DB.Exec("drop table follow")
	DB.Exec("drop table users")
}
