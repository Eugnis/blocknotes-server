package main

import (
	"blocknotes_server/pkg/mongo"
	"blocknotes_server/pkg/server"
	"log"
)

func main() {

	mongouri := "mongodb://blocknotes-mongodb:27017"
	mongouri = "127.0.0.1:27017"
	ms, err := mongo.NewSession(mongouri)
	if err != nil {
		log.Fatalln("unable to connect to mongodb" + mongouri)
	}
	defer ms.Close()

	// scheduler.Start()
	mongoDbName := "blocknotes_server"
	ns := mongo.NewNoteService(ms.Copy(), mongoDbName, "note")

	s := server.NewServer(ns)

	s.Start()
}
