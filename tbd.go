package main

import (
	"tbd/service"
)

func main() {
	router := service.InitRouter()
	router.Run(":5005")
}
