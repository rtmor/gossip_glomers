package main

import (
	"log"

	"gglomers/internal/broadcast"
)

func main() {
	s := broadcast.NewBroadcastService()
	log.Panic(s.Run())
}
