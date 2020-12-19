package main

import (
	_ "tbd/model"
	"tbd/router"
	_ "tbd/router"
)

func main() {
	// testData()
	router := router.InitRouter()
	router.Run(":5005")
	// fmt.Println(model.QueryReplyWithID(5))
	// fmt.Println(model.QueryAllCommentsWithContentID(2))
}
