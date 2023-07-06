package main

import (
	"log"

	"gglomers/internal/uniqueids"
)

func main() {
	node := uniqueids.NewUniqueIDService()
	log.Panic(node.Run())
}
