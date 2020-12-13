package model

func init() {
	Connect()
	CreateUserTableIfNotExists()
	CreateFollowTableIfNotExists()
	CreateContentTableIfNotExists()
	CreateImageTableIfNotExists()
}
