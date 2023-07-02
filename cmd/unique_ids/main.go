package main

import (
	"log"

	gg "gglomers"
)

func main() {
	node := gg.NewUniqueIDService()
	log.Panic(node.Run())
}
