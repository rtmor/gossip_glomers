package main

import (
	"log"

	gg "gglomers"
)

func main() {
	s := gg.NewBroadcastService()
	log.Panic(s.Run())
}
